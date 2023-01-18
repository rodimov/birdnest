package delivery

import (
	"birdnest/internal/drones/delivery/dto"
	"birdnest/internal/drones/model"
	"math"
)

func getDistance(x float64, y float64) int64 {
	const XC float64 = 250000
	const YC float64 = 250000

	return int64(math.Round(math.Sqrt(math.Pow(x-XC, 2) + math.Pow(y-YC, 2))))
}

func ToDroneDto(from *model.Drone) *dto.DroneDto {
	target := &dto.DroneDto{
		ID:        from.ID,
		PositionX: from.PositionX,
		PositionY: from.PositionY,
		LastSeen:  from.LastSeen,
		Distance:  getDistance(from.PositionX, from.PositionY),
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
			Distance:  getDistance(drone.PositionX, drone.PositionY),
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
