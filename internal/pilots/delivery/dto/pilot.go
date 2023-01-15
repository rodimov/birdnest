package dto

import (
	"birdnest/internal/drones/delivery/dto"
)

type PilotDto struct {
	ID               string        `json:"id"`
	FirstName        string        `json:"first_name"`
	LastName         string        `json:"last_name"`
	Phone            string        `json:"phone"`
	Email            string        `json:"email"`
	RegistrationTime string        `json:"registration_time"`
	Drone            *dto.DroneDto `json:"drone"`
}
