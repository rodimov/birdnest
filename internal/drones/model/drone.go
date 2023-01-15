package model

import (
	"database/sql"
	"time"
)

type Drone struct {
	ID        string       `json:"id"`
	PositionX float64      `json:"position_x"`
	PositionY float64      `json:"position_y"`
	LastSeen  time.Time    `json:"last_seen"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
