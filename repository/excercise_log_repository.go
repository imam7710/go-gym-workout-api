package repository

import (
	"gym-app/domain"

	"gorm.io/gorm"
)

type exerciseLogRepository struct {
	db *gorm.DB
}

func NewExerciseLogRepository(db *gorm.DB) domain.ExerciseLogRepository {
	return &exerciseLogRepository{db: db}
}

func (r *exerciseLogRepository) Create(log *domain.ExerciseLog) error {
	return r.db.Create(log).Error
}

func (r *exerciseLogRepository) GetByExerciseID(exerciseID int) ([]domain.ExerciseLog, error) {
	var logs []domain.ExerciseLog
	err := r.db.Where("exercise_id = ?", exerciseID).Order("created_at desc").Find(&logs).Error
	return logs, err
}
