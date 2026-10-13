package usecase

import (
	"errors"
	"gym-app/domain"

	"github.com/go-playground/validator/v10"
)

type workoutUsecase struct {
	workoutRepo domain.WorkoutRepository
	validate    *validator.Validate
}

func NewWorkoutUsecase(workoutRepo domain.WorkoutRepository, validate *validator.Validate) domain.WorkoutUsecase {
	return &workoutUsecase{
		workoutRepo: workoutRepo,
		validate:    validate,
	}
}

func (u *workoutUsecase) CreateWorkout(userID int, req domain.WorkoutRequest) (*domain.Workout, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, errors.New("data request tidak valid: " + err.Error())
	}

	workout := &domain.Workout{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	err := u.workoutRepo.Create(workout)
	if err != nil {
		return nil, errors.New("gagal menyimpan data workout")
	}

	return workout, nil
}

func (u *workoutUsecase) GetWorkoutsByUserID(userID int) ([]domain.Workout, error) {
	return u.workoutRepo.GetByUserID(userID)
}

// Tambahan: GET Workout by ID
func (u *workoutUsecase) GetWorkoutByID(id int, userID int) (*domain.Workout, error) {
	workout, err := u.workoutRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return nil, errors.New("data tidak ditemukan")
	}
	return workout, nil
}

// Tambahan: UPDATE (PUT) Workout
func (u *workoutUsecase) UpdateWorkout(id int, userID int, req domain.WorkoutRequest) (*domain.Workout, error) {
	// 1. Validasi request
	if err := u.validate.Struct(req); err != nil {
		return nil, errors.New("data request tidak valid: " + err.Error())
	}

	// 2. Cari data lama (sekaligus memastikan kepemilikan user)
	workout, err := u.workoutRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return nil, errors.New("data tidak ditemukan")
	}

	// 3. Timpa dengan data baru
	workout.Name = req.Name
	workout.Description = req.Description

	// 4. Simpan ke database
	err = u.workoutRepo.Update(workout)
	if err != nil {
		return nil, errors.New("gagal memperbarui data workout")
	}

	return workout, nil
}

// Perubahan: DELETE Workout dengan validasi 404 dan return data terhapus
func (u *workoutUsecase) DeleteWorkout(id int, userID int) (*domain.Workout, error) {
	// 1. Cari datanya dulu. Jika tidak ada, kembalikan pesan "data tidak ditemukan" (untuk 404)
	workout, err := u.workoutRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return nil, errors.New("data tidak ditemukan")
	}

	// 2. Jika ketemu, lakukan penghapusan
	err = u.workoutRepo.Delete(workout)
	if err != nil {
		return nil, errors.New("error server saat menghapus data") // (untuk 500)
	}

	// 3. Kembalikan data yang baru saja dihapus
	return workout, nil
}
