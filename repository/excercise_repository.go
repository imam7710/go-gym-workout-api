package repository

import (
	"gym-app/domain"

	"gorm.io/gorm"
)

type exerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) domain.ExerciseRepository {
	return &exerciseRepository{db: db}
}

func (r *exerciseRepository) Create(exercise *domain.Exercise) error {
	return r.db.Create(exercise).Error
}

func (r *exerciseRepository) GetByWorkoutID(workoutID int) ([]domain.Exercise, error) {
	var exercises []domain.Exercise
	err := r.db.Where("workout_id = ?", workoutID).Find(&exercises).Error
	return exercises, err
}

func (r *exerciseRepository) GetByID(id int) (*domain.Exercise, error) {
	var exercise domain.Exercise
	err := r.db.Where("id = ?", id).First(&exercise).Error
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (r *exerciseRepository) Delete(exercise *domain.Exercise) error {
	return r.db.Delete(exercise).Error
}
