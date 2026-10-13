package usecase

import (
	"errors"
	"gym-app/domain"

	"github.com/go-playground/validator/v10"
)

type exerciseUsecase struct {
	exerciseRepo domain.ExerciseRepository
	workoutRepo  domain.WorkoutRepository
	validate     *validator.Validate
}

func NewExerciseUsecase(exerciseRepo domain.ExerciseRepository, workoutRepo domain.WorkoutRepository, validate *validator.Validate) domain.ExerciseUsecase {
	return &exerciseUsecase{
		exerciseRepo: exerciseRepo,
		workoutRepo:  workoutRepo,
		validate:     validate,
	}
}

// 1. FUNGSI CREATE
func (u *exerciseUsecase) CreateExercise(userID int, req domain.ExerciseRequest) (*domain.Exercise, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, errors.New("data request tidak valid: " + err.Error())
	}

	// Verifikasi kepemilikan Workout
	_, err := u.workoutRepo.GetByIDAndUserID(req.WorkoutID, userID)
	if err != nil {
		return nil, errors.New("unauthorized: workout tidak ditemukan atau bukan milik anda")
	}

	exercise := &domain.Exercise{
		WorkoutID:   req.WorkoutID,
		Name:        req.Name,
		Description: req.Description,
	}

	err = u.exerciseRepo.Create(exercise)
	if err != nil {
		return nil, errors.New("gagal menyimpan data exercise")
	}

	return exercise, nil
}

// 2. FUNGSI GET
func (u *exerciseUsecase) GetExercisesByWorkoutID(workoutID int) ([]domain.Exercise, error) {
	return u.exerciseRepo.GetByWorkoutID(workoutID)
}

// 3. FUNGSI DELETE
func (u *exerciseUsecase) DeleteExercise(id int, userID int) (*domain.Exercise, error) {
	// Cari data exercise-nya (404)
	exercise, err := u.exerciseRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("data exercise tidak ditemukan")
	}

	// Verifikasi kepemilikan: Apakah workout dari exercise ini milik user yang login?
	_, err = u.workoutRepo.GetByIDAndUserID(exercise.WorkoutID, userID)
	if err != nil {
		return nil, errors.New("unauthorized: tidak dapat menghapus exercise pada workout orang lain")
	}

	// Hapus (500)
	err = u.exerciseRepo.Delete(exercise)
	if err != nil {
		return nil, errors.New("error server saat menghapus data")
	}

	return exercise, nil
}
