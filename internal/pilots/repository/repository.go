package repository

import (
	pilotsPkg "birdnest/internal/pilots"
	"birdnest/internal/pilots/model"
	"birdnest/logger"
	"context"
	"errors"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewPilotsRepository(dbInstance *gorm.DB) pilotsPkg.Repository {
	logger.AppLogger.Info("Pilots repository was created.")
	return &repository{db: dbInstance}
}

func (r *repository) Create(context context.Context, pilot *model.Pilot) (*model.Pilot, error) {
	result := r.db.Create(pilot)

	if result.Error != nil {
		return nil, result.Error
	} else {
		return pilot, nil
	}
}

func (r *repository) GetByDroneId(context context.Context, droneId string) (*model.Pilot, error) {
	var pilot model.Pilot
	result := r.db.Where("drone_id = ?", droneId).First(&pilot)

	if result.RowsAffected == 0 {
		return nil, errors.New("pilot not found by drone_id " + droneId)
	}

	return &pilot, nil
}

func (r *repository) IsPilotExists(context context.Context, droneId string) bool {
	var exists bool

	err := r.db.Model(model.Pilot{}).
		Select("count(*) > 0").
		Where("drone_id = ?", droneId).
		Find(&exists).
		Error

	if err != nil {
		return false
	}

	return exists
}
