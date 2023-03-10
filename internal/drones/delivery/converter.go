package delivery

import (
	"birdnest/internal/drones/delivery/dto"
	"birdnest/internal/drones/model"
	"birdnest/pkg"
)

func ToDroneDto(from *model.Drone) *dto.DroneDto {
	target := &dto.DroneDto{
		ID:        from.ID,
		PositionX: from.PositionX,
		PositionY: from.PositionY,
		LastSeen:  from.LastSeen,
		Distance:  pkg.GetDroneDistance(from.PositionX, from.PositionY),
	}

	return target
}

func ToDronesDto(from []*model.Drone) []*dto.DroneDto {
	var target []*dto.DroneDto

	for _, drone := range from {
		droneDto := &dto.DroneDto{
			ID:        drone.ID,
			PositionX: drone.PositionX,
			PositionY: drone.PositionY,
			LastSeen:  drone.LastSeen,
			Distance:  pkg.GetDroneDistance(drone.PositionX, drone.PositionY),
		}

		target = append(target, droneDto)
	}

	return target
}

func ToDroneModel(from *dto.DroneDto) *model.Drone {
	target := &model.Drone{
		ID:        from.ID,
		PositionX: from.PositionX,
		PositionY: from.PositionY,
		LastSeen:  from.LastSeen,
	}

	return target
}
