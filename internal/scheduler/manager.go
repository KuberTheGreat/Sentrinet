package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/KuberTheGreat/Sentrinet/internal/scan"
	"github.com/jmoiron/sqlx"
)

type JobRow struct{
	ID int64 `db: "id"`
	Target string `db:"target"`
	StartPort int `db:"start_port"`
	EndPort int `db:"end_port"`
	IntervalSeconds int `db:"interval_seconds"`
	Active int `db:"active"`
	CreatedAt time.Time `db:"created_at"`
}

type Manager struct{
	db *sqlx.DB
	ctx context.Context
	cancel context.CancelFunc
	mu sync.Mutex
	runners map[int64]*jobRunner
	wg sync.WaitGroup
}

type jobRunner struct{
	cancel context.CancelFunc
	running int32
	jobRow JobRow
}

func NewManager(parentCtx context.Context, db *sqlx.DB) *Manager{
	ctx, cancel := context.WithCancel(parentCtx)
	return &Manager{
		db: db,
		ctx: ctx,
		cancel: cancel,
		runners: make(map[int64]*jobRunner),
	}
}

func (m *Manager) LoadAndStartAll() error{
	rows := []JobRow{}
	if err := m.db.Select(&rows, "SELECT * FROM jobs WHERE active = 1"); err != nil{
		if err == sql.ErrNoRows{
			return nil
		}
		return err
	}

	for _, jr := range rows{
		if err := m.startRunner(jr); err != nil{
			fmt.Printf("[Scheduler] failed to start job %d: %v\n", jr.ID, err)
		}
	}
	return nil
}

func (m *Manager) CreateJob(target string, startPort, endPort int, interval time.Duration, active bool) (int64, error){
	intervalSec := int(interval.Seconds())
	activeInt := 0
	if active{
		activeInt = 1
	}

	res, err := m.db.Exec(
		`INSERT INTO jobs (target, start_port, end_port, interval_seconds, active)
		VALUES (?, ?, ?, ?, ?)`,
		target, startPort, endPort, intervalSec, activeInt,
	)
	if err != nil{
		return 0, err
	}

	id, _ := res.LastInsertId()

	if active{
		jr := JobRow{
			ID: id,
			Target: target,
			StartPort: startPort,
			EndPort: endPort,
			IntervalSeconds: intervalSec,
			Active: activeInt,
		}
		if err := m.startRunner(jr); err != nil{
			fmt.Printf("[Scheduler] created job %d but failed to start runner: %v\n", id, err)
		}
	}

	return id, nil
}

func (m *Manager) startRunner(jr JobRow) error{
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.runners[jr.ID]; exists {
		return fmt.Errorf("job %d already running", jr.ID)
	}

	ctx, cancel := context.WithCancel(m.ctx)
	runner := &jobRunner{
		cancel: cancel,
		jobRow: jr,
	}

	m.runners[jr.ID] = runner
	m.wg.Add(1)

	go func(){
		defer m.wg.Done()
		ticker := time.NewTicker(time.Duration(jr.IntervalSeconds) * time.Second)
		defer ticker.Stop()

		if err := m.executeOnce(ctx, jr, runner); err != nil{
			fmt.Printf("[Scheduler] initial run error for job %d: %v\n", jr.ID, err)
		}

		for{
			select{
			case <-ctx.Done():
				fmt.Printf("[Scheduler] job %d stopped\n", jr.ID)
				return
			case <-ticker.C:
				if !atomic.CompareAndSwapInt32(&runner.running, 0, 1){
					fmt.Printf("[Scheduler] job %d previous run still active, skipping tick\n", jr.ID)
					continue
				}

				go func(){
					defer atomic.StoreInt32(&runner.running, 0)
					if err := m.executeScanAndSave(jr); err != nil{
						fmt.Printf("[Scheduler] job %d run error: %v\n", jr.ID, err)
					}
				}()
			}
		}
	}()

	return nil
}

func (m *Manager) executeOnce(ctx context.Context, jr JobRow, runner *jobRunner) error{
	if !atomic.CompareAndSwapInt32(&runner.running, 0, 1){
		return nil
	}
	defer atomic.StoreInt32(&runner.running, 0)

	select{
	case <-ctx.Done():
		return fmt.Errorf("job %d cancelled before start", jr.ID)
	default:
	}

	if err := m.executeScanAndSave(jr); err != nil{
		return err
	}

	return nil
}

func(m *Manager) executeScanAndSave(jr JobRow) error{
	fmt.Printf("[Scheduler] Running recurring scan for %s (job %d)\n", jr.Target, jr.ID)
	results := scan.ScanRange(jr.Target, jr.StartPort, jr.EndPort)
	tx, err := m.db.Beginx()
	if err != nil{
		for _, r := range results{
			if _, e := m.db.NamedExec(
				`INSERT INTO scans (target, port, is_open, duration_ms) VALUES (:target, :port, :is_open, :duration_ms)`,
				map[string]interface{}{
					"target": jr.Target,
					"port": r.Port,
					"is_open": r.IsOpen,
					"duration_ms": r.Duration,
				},
			); e != nil {
				fmt.Printf("[Scheduler] recurring insert error: %v\n", e)
			}
		}
		return err
	}

	for _, r := range results{
		if _, e := tx.NamedExec(
			`INSERT INTO scans (target, port, is_open, duration_ms) VALUES (:target, :port, :is_open, :duration_ms)`,
			map[string]interface{}{
				"target": jr.Target,
				"port": r.Port,
				"is_open": r.IsOpen,
				"duration_ms": r.Duration,
			},
		); e != nil {
			fmt.Printf("[Scheduler] recurring insert error(tx): %v\n", e)
		}
	}

	if err := tx.Commit(); err != nil{
		fmt.Printf("[Scheduler] tx commit error: %v\n", err)
		return err
	}

	return nil
}

func (m *Manager) StopJob(id int64) error{
	m.mu.Lock()
	runner, ok := m.runners[id]
	m.mu.Unlock()

	if ok {
		runner.cancel()
		m.mu.Lock()
		delete(m.runners, id)
		m.mu.Unlock()
	}

	_, err := m.db.Exec("UPDATE jobs SET active = 0 WHERE id = ?", id)
	return err
}

func (m *Manager) StartJobByID(id int64) error{
	var jr JobRow
	if err := m.db.Get(&jr, "SELECT * FROM jobs WHERE id = ?", id); err != nil{
		return err
	}
	if jr.Active == 1{
		return fmt.Errorf("job already active")
	}

	if _, err := m.db.Exec("UPDATE jobs SET active = 1 WHERE id  = ?", id); err != nil{
		return err
	}

	jr.Active = 1
	return m.startRunner(jr)
}

func (m *Manager) DeleteJob(id int64) error{
	_ = m.StopJob(id)
	_, err := m.db.Exec("DELETE FROM jobs WHERE id = ?", id)
	return err
}

func (m *Manager) ListJobs() ([]JobRow, error){
	rows := []JobRow{}
	if err := m.db.Select(&rows, "SELECT * FROM jobs ORDER BY created_at DESC"); err != nil{
		return nil, err
	}
	return rows, nil
}

func (m *Manager) StopAll(){
	m.cancel()
	m.mu.Lock()
	for _, r := range m.runners{
		r.cancel()
	}
	m.mu.Unlock()
	m.wg.Wait()
}