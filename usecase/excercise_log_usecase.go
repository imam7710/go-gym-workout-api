package usecase

import (
	"errors"
	"gym-app/domain"

	"github.com/go-playground/validator/v10"
)

type exerciseLogUsecase struct {
	logRepo      domain.ExerciseLogRepository
	exerciseRepo domain.ExerciseRepository
	workoutRepo  domain.WorkoutRepository
	validate     *validator.Validate
}

func NewExerciseLogUsecase(logRepo domain.ExerciseLogRepository, exerciseRepo domain.ExerciseRepository, workoutRepo domain.WorkoutRepository, validate *validator.Validate) domain.ExerciseLogUsecase {
	return &exerciseLogUsecase{
		logRepo:      logRepo,
		exerciseRepo: exerciseRepo,
		workoutRepo:  workoutRepo,
		validate:     validate,
	}
}

func (u *exerciseLogUsecase) CreateLog(userID int, req domain.ExerciseLogRequest) (*domain.ExerciseLog, error) {
	// 1. Validasi Input
	if err := u.validate.Struct(req); err != nil {
		return nil, errors.New("data request tidak valid: " + err.Error())
	}

	// 2. Cari Exercise-nya
	exercise, err := u.exerciseRepo.GetByID(req.ExerciseID)
	if err != nil {
		return nil, errors.New("exercise tidak ditemukan")
	}

	// 3. Verifikasi Kepemilikan (Apakah Workout dari Exercise ini milik User yang login?)
	_, err = u.workoutRepo.GetByIDAndUserID(exercise.WorkoutID, userID)
	if err != nil {
		return nil, errors.New("unauthorized: anda tidak memiliki akses ke exercise ini")
	}

	// 4. Simpan Log (Penyesuaian dengan struktur DB)
	log := &domain.ExerciseLog{
		ExerciseID: req.ExerciseID,
		UserID:     userID,
		SetCount:   req.SetCount,
		RepCount:   req.RepCount,
		Weight:     req.Weight,
	}

	err = u.logRepo.Create(log)
	if err != nil {
		return nil, errors.New("gagal menyimpan data log")
	}

	return log, nil
}

func (u *exerciseLogUsecase) GetLogsByExerciseID(exerciseID int, userID int) ([]domain.ExerciseLog, error) {
	// 1. Cari Exercise-nya
	exercise, err := u.exerciseRepo.GetByID(exerciseID)
	if err != nil {
		return nil, errors.New("exercise tidak ditemukan")
	}

	// 2. Verifikasi Kepemilikan (Mencegah user mengintip log orang lain)
	_, err = u.workoutRepo.GetByIDAndUserID(exercise.WorkoutID, userID)
	if err != nil {
		return nil, errors.New("unauthorized: anda tidak memiliki akses ke exercise ini")
	}

	// 3. Ambil data logs
	return u.logRepo.GetByExerciseID(exerciseID)
}
