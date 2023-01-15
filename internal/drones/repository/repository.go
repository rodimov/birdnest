package repository

import (
	dr "birdnest/internal/drones"
	"birdnest/internal/drones/model"
	"birdnest/logger"
	"birdnest/pkg"
	"context"
	"gorm.io/gorm"
	"time"
)

type repository struct {
	db *gorm.DB
}

func DronesRepository(dbInstance *gorm.DB) dr.Repository {
	logger.AppLogger.Info("Drone repository was created.")
	return &repository{db: dbInstance}
}

func (r *repository) Create(context context.Context, drone *model.Drone) (*model.Drone, error) {
	result := r.db.Create(drone)

	if result.Error != nil {
		return nil, result.Error
	} else {
		return drone, nil
	}
}

func (r *repository) Update(context context.Context, drone *model.Drone) (*model.Drone, error) {
	result := r.db.Save(drone)

	if result.Error != nil {
		return nil, result.Error
	}

	return drone, nil
}

func (r *repository) GetAll(context context.Context) ([]*model.Drone, error) {
	var drones []model.Drone
	result := r.db.Where("deleted_at is null").Find(&drones)

	if result.Error != nil {
		return nil, result.Error
	}

	dronesPtr := make([]*model.Drone, len(drones))

	for i, _ := range drones {
		dronesPtr[i] = &drones[i]
	}

	return dronesPtr, nil
}

func (r *repository) IsDroneExists(context context.Context, id string) (bool, error) {
	var exists bool

	err := r.db.Model(model.Drone{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *repository) IsDroneRemoved(context context.Context, id string) (bool, error) {
	var isNotRemoved bool

	err := r.db.Model(model.Drone{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Where("deleted_at is null").
		Find(&isNotRemoved).
		Error

	if err != nil {
		return true, err
	}

	return !isNotRemoved, nil
}

func (r *repository) GetById(context context.Context, id string) (*model.Drone, error) {
	drone := model.Drone{ID: id}

	result := r.db.Take(&drone)

	if result.Error != nil {
		return nil, result.Error
	} else {
		return &drone, nil
	}
}

func (r *repository) DeleteById(context context.Context, id string) (*model.Drone, error) {
	drone := model.Drone{ID: id}
	result := r.db.Take(&drone)

	if result.Error != nil {
		return nil, result.Error
	}

	drone.DeletedAt = pkg.TimeToNullTime(time.Now())

	return r.Update(context, &drone)
}

func (r *repository) DeletePermanentlyById(context context.Context, id string) (*model.Drone, error) {
	drone := model.Drone{ID: id}
	result := r.db.Take(&drone)

	if result.Error != nil {
		return nil, result.Error
	}

	result = r.db.Delete(drone)

	if result.Error != nil {
		return nil, result.Error
	}

	return &drone, nil
}

func (r *repository) GetAllDeleted(context context.Context) ([]*model.Drone, error) {
	var drones []model.Drone
	result := r.db.Where("deleted_at is not null").Find(&drones)

	if result.Error != nil {
		return nil, result.Error
	}

	dronesPtr := make([]*model.Drone, len(drones))

	for i, _ := range drones {
		dronesPtr[i] = &drones[i]
	}

	return dronesPtr, nil
}
