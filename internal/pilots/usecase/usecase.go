package usecase

import (
	dronesPkg "birdnest/internal/drones"
	"birdnest/internal/pilots"
	"birdnest/internal/pilots/model"
	"context"
)

type usecase struct {
	Repository pilots.Repository
	DronesUC   dronesPkg.UseCase
}

func NewUseCase(repository pilots.Repository, dronesUC dronesPkg.UseCase) pilots.UseCase {
	return &usecase{Repository: repository, DronesUC: dronesUC}
}

func (u *usecase) Create(context context.Context, pilot *model.Pilot) (*model.Pilot, error) {
	return u.Repository.Create(context, pilot)
}

func (u *usecase) GetByDroneId(context context.Context, droneId string) (*model.Pilot, error) {
	pilot, err := u.Repository.GetByDroneId(context, droneId)

	if err != nil {
		return nil, err
	}

	drone, err := u.DronesUC.GetById(context, droneId)

	if err != nil {
		return nil, err
	}

	pilot.Drone = drone

	return pilot, nil
}

func (u *usecase) IsExistsByDroneId(context context.Context, droneId string) bool {
	return u.Repository.IsPilotExists(context, droneId)
}

func (u *usecase) GetAll(context context.Context) ([]*model.Pilot, error) {
	drones, err := u.DronesUC.GetAll(context)

	if err != nil {
		return nil, err
	}

	var pilotsSlice []*model.Pilot

	for _, drone := range drones {
		pilot, err := u.GetByDroneId(context, drone.ID)

		if err != nil {
			return nil, err
		}

		pilotsSlice = append(pilotsSlice, pilot)
	}

	return pilotsSlice, nil
}
