package repository

import (
	"gym-app/domain"

	"gorm.io/gorm"
)

type workoutRepository struct {
	db *gorm.DB
}

func NewWorkoutRepository(db *gorm.DB) domain.WorkoutRepository {
	return &workoutRepository{db: db}
}

func (r *workoutRepository) Create(workout *domain.Workout) error {
	return r.db.Create(workout).Error
}

func (r *workoutRepository) GetByUserID(userID int) ([]domain.Workout, error) {
	var workouts []domain.Workout
	// Tambahan: Gunakan Preload agar array exercises ikut terbawa
	err := r.db.Preload("Exercises").Where("user_id = ?", userID).Find(&workouts).Error
	return workouts, err
}

// Tambahan: Fungsi untuk mencari 1 workout spesifik (dipakai untuk GET by ID, PUT, dan DELETE)
func (r *workoutRepository) GetByIDAndUserID(id int, userID int) (*domain.Workout, error) {
	var workout domain.Workout
	err := r.db.Preload("Exercises").Where("id = ? AND user_id = ?", id, userID).First(&workout).Error
	if err != nil {
		return nil, err
	}
	return &workout, nil
}

// Tambahan: Fungsi untuk menyimpan perubahan data (PUT)
func (r *workoutRepository) Update(workout *domain.Workout) error {
	return r.db.Save(workout).Error
}

// Perbaikan: Hapus data dengan struct yang sudah ditemukan sebelumnya
func (r *workoutRepository) Delete(workout *domain.Workout) error {
	return r.db.Delete(workout).Error
}
