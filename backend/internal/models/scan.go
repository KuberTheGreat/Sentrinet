package models

import "time"

type ScanResult struct {
	ID int64 `db:"id" json:"id"`
	Target string `db:"target" json:"target"`
	Port int `db:"port" json:"port"`
	IsOpen bool `db:"is_open" json:"is_open"`
	Duration int64 `db:"duration_ms" json:"duration_ms"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UserId int64 `db:"user_id" json:"user_id"`
}