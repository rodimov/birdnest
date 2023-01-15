package usecase

import (
	dr "birdnest/internal/drones"
	"birdnest/internal/drones/model"
	"context"
)

type usecase struct {
	Repository dr.Repository
}

func DronesUseCase(repository dr.Repository) dr.UseCase {
	return &usecase{Repository: repository}
}

func (u *usecase) Create(context context.Context, drone *model.Drone) (*model.Drone, error) {
	isDroneExists, err := u.Repository.IsDroneExists(context, drone.ID)

	if err != nil {
		return nil, err
	}

	if isDroneExists {
		_, err = u.Repository.Update(context, drone)

		if err != nil {
			return nil, err
		}
	} else {
		_, err = u.Repository.Create(context, drone)

		if err != nil {
			return nil, err
		}
	}

	return drone, nil
}

func (u *usecase) GetAll(context context.Context) ([]*model.Drone, error) {
	return u.Repository.GetAll(context)
}

func (u *usecase) GetById(context context.Context, id string) (*model.Drone, error) {
	return u.Repository.GetById(context, id)
}

func (u *usecase) IsRemoved(context context.Context, id string) (bool, error) {
	return u.Repository.IsDroneRemoved(context, id)
}

func (u *usecase) DeleteById(context context.Context, id string) (*model.Drone, error) {
	return u.Repository.DeleteById(context, id)
}

func (u *usecase) DeletePermanentlyById(context context.Context, id string) (*model.Drone, error) {
	return u.Repository.DeletePermanentlyById(context, id)
}

func (u *usecase) GetAllDeleted(context context.Context) ([]*model.Drone, error) {
	return u.Repository.GetAllDeleted(context)
}
