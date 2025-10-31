package models

import "time"

type Notification struct {
    ID        int       `db:"id" json:"id"`
    UserID    int       `db:"user_id" json:"user_id"`
    ScanID    int       `db:"scan_id" json:"scan_id"`
    Type      string    `db:"type" json:"type"`
    Message   string    `db:"message" json:"message"`
    Read      bool      `db:"read" json:"read"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}
