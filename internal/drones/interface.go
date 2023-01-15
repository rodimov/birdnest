package drones

import (
	"birdnest/internal/drones/model"
	"context"
)

type UseCase interface {
	Create(context context.Context, drone *model.Drone) (*model.Drone, error)
	GetAll(context context.Context) ([]*model.Drone, error)
	GetById(context context.Context, id string) (*model.Drone, error)
	IsRemoved(context context.Context, id string) (bool, error)
	DeleteById(context context.Context, id string) (*model.Drone, error)
}

type Repository interface {
	Create(context context.Context, drone *model.Drone) (*model.Drone, error)
	Update(context context.Context, drone *model.Drone) (*model.Drone, error)
	IsDroneExists(context context.Context, id string) (bool, error)
	IsDroneRemoved(context context.Context, id string) (bool, error)
	GetAll(context context.Context) ([]*model.Drone, error)
	GetById(context context.Context, id string) (*model.Drone, error)
	DeleteById(context context.Context, id string) (*model.Drone, error)
}
