package pilots

import (
	"birdnest/internal/pilots/model"
	"context"
)

type UseCase interface {
	Create(context context.Context, pilot *model.Pilot) (*model.Pilot, error)
	GetByDroneId(context context.Context, droneId string) (*model.Pilot, error)
	GetAll(context context.Context) ([]*model.Pilot, error)
	IsExistsByDroneId(context context.Context, droneId string) bool
}

type Repository interface {
	Create(context context.Context, pilot *model.Pilot) (*model.Pilot, error)
	GetByDroneId(context context.Context, droneId string) (*model.Pilot, error)
	IsPilotExists(context context.Context, droneId string) bool
}
