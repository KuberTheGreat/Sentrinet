package scheduler

import (
	"fmt"
	"time"

	"github.com/KuberTheGreat/Sentrinet/internal/scan"
	"github.com/jmoiron/sqlx"
)

type Job struct{
	Target string
	StartPort int
	EndPort int
	Interval time.Duration
}

func StartJob(db *sqlx.DB, job Job, userId int64){
	ticker := time.NewTicker(job.Interval)
	go func(){
		for{
			<-ticker.C
			fmt.Printf("[Scheduler] Running recurring scan for %s\n", job.Target)

			results := scan.ScanRange(job.Target, job.StartPort, job.EndPort)
			for _, r := range results{
				_, err := db.NamedExec(
					`INSERT INTO scans (target, port, is_open, duration_ms, user_id)
					VALUES (:target, :port, :is_open, :duration_ms, :user_id)`,
					map[string]interface{}{
						"target": job.Target,
						"port": r.Port,
						"is_open": r.IsOpen,
						"duration_ms": r.Duration,
						"user_id": userId,
					},
				)
				if err != nil{
					fmt.Println("Recurring insert error: ", err)
				}
			}
		}
	}()
}