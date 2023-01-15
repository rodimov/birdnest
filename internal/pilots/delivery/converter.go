package delivery

import (
	dronesDelivery "birdnest/internal/drones/delivery"
	"birdnest/internal/pilots/delivery/dto"
	"birdnest/internal/pilots/model"
	"birdnest/pkg"
	"time"
)

func ToPilotDto(from *model.Pilot) *dto.PilotDto {
	target := &dto.PilotDto{
		ID:               from.ID,
		FirstName:        pkg.NullStringToString(from.FirstName),
		LastName:         pkg.NullStringToString(from.LastName),
		Phone:            pkg.NullStringToString(from.Phone),
		Email:            pkg.NullStringToString(from.Email),
		RegistrationTime: pkg.NullTimeToString(from.RegistrationTime),
		Drone:            dronesDelivery.ToDroneDto(from.Drone),
	}

	return target
}

func ToPilotsDto(from []*model.Pilot) []*dto.PilotDto {
	var target []*dto.PilotDto

	for _, pilot := range from {
		pilotDto := &dto.PilotDto{
			ID:               pilot.ID,
			FirstName:        pkg.NullStringToString(pilot.FirstName),
			LastName:         pkg.NullStringToString(pilot.LastName),
			Phone:            pkg.NullStringToString(pilot.Phone),
			Email:            pkg.NullStringToString(pilot.Email),
			RegistrationTime: pkg.NullTimeToString(pilot.RegistrationTime),
			Drone:            dronesDelivery.ToDroneDto(pilot.Drone),
		}

		target = append(target, pilotDto)
	}

	return target
}

func ToPilot(from *dto.PilotDto) *model.Pilot {
	registrationTime, _ := time.Parse(time.RFC3339, from.RegistrationTime)
	
	target := &model.Pilot{
		ID:               from.ID,
		FirstName:        pkg.StringToNullString(from.FirstName),
		LastName:         pkg.StringToNullString(from.LastName),
		Phone:            pkg.StringToNullString(from.Phone),
		Email:            pkg.StringToNullString(from.Email),
		RegistrationTime: pkg.TimeToNullTime(registrationTime),
		Drone:            dronesDelivery.ToDroneModel(from.Drone),
	}

	return target
}
