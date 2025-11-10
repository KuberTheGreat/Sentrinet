package db

import (
	"fmt"
	"log"
	"time"

	"github.com/KuberTheGreat/Sentrinet/internal/metrics"
	"github.com/jmoiron/sqlx"
)

func StartCleanupScheduler(database *sqlx.DB){
	go func() {
		ticker := time.NewTicker(60 * time.Minute)
		defer ticker.Stop()

		for{
			if err := removeClosedPortsOlderThan(database, time.Hour); err!= nil{
				log.Println("[Cleanup Error]: ", err)
			}
			<-ticker.C
		}
	}()
}

func removeClosedPortsOlderThan(db *sqlx.DB, olderThan time.Duration) error{
	start := time.Now()
	cutoff := time.Now().Add(-olderThan)

	res, err := db.Exec("DELETE FROM scans WHERE is_open = 0 AND created_at < ?", cutoff)
	if err != nil{
		return err
	}

	count, _ := res.RowsAffected()
	duration := time.Since(start).Milliseconds()

	_, logErr := db.Exec(`INSERT INTO cleanup_logs (deleted_count, run_time_ms) VALUEs (?, ?)`,count, duration)
	
	if logErr != nil{
		fmt.Println("[Cleanup Log Error]: ", logErr)
	}
	
	if count > 0{
		fmt.Printf("[Cleanup] Deleted %d closed ports older than %v\n", count, olderThan)
		metrics.CleanupDeleted.Add(float64(count))
	}
	return nil
} 