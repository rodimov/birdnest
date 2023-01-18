package dto

import "time"

type DroneDto struct {
	ID        string    `json:"id"`
	PositionX float64   `json:"position_x"`
	PositionY float64   `json:"position_y"`
	LastSeen  time.Time `json:"last_seen"`
	Distance  int64     `json:"distance"`
}
