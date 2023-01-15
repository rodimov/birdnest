package model

import (
	"birdnest/internal/drones/model"
	"database/sql"
)

type Pilot struct {
	ID               string         `json:"id"`
	FirstName        sql.NullString `json:"first_name"`
	LastName         sql.NullString `json:"last_name"`
	Phone            sql.NullString `json:"phone"`
	Email            sql.NullString `json:"email"`
	RegistrationTime sql.NullTime   `json:"registration_time"`
	DroneID          string         `json:"drone_id"`
	Drone            *model.Drone   `json:"drone"`
}
