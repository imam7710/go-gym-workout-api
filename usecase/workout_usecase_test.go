package usecase_test

import (
	"errors"
	"gym-app/domain"
	"gym-app/usecase"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==========================================
// 1. SETUP MOCK REPOSITORY
// ==========================================
type MockWorkoutRepository struct {
	mock.Mock
}

func (m *MockWorkoutRepository) Create(workout *domain.Workout) error {
	args := m.Called(workout)
	return args.Error(0)
}

func (m *MockWorkoutRepository) GetByUserID(userID int) ([]domain.Workout, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Workout), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWorkoutRepository) GetByIDAndUserID(id int, userID int) (*domain.Workout, error) {
	args := m.Called(id, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Workout), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWorkoutRepository) Update(workout *domain.Workout) error {
	args := m.Called(workout)
	return args.Error(0)
}

func (m *MockWorkoutRepository) Delete(workout *domain.Workout) error {
	args := m.Called(workout)
	return args.Error(0)
}

// ==========================================
// 2. UNIT TEST UNTUK FUNGSI GET BY ID
// ==========================================
func TestGetWorkoutByID_Success(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	validator := validator.New()
	workoutUsecase := usecase.NewWorkoutUsecase(mockRepo, validator)

	mockWorkout := &domain.Workout{ID: 1, UserID: 2, Name: "Leg Day"}
	mockRepo.On("GetByIDAndUserID", 1, 2).Return(mockWorkout, nil)

	result, err := workoutUsecase.GetWorkoutByID(1, 2)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Leg Day", result.Name)
}

func TestGetWorkoutByID_NotFound(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	validator := validator.New()
	workoutUsecase := usecase.NewWorkoutUsecase(mockRepo, validator)

	mockRepo.On("GetByIDAndUserID", 99, 2).Return(nil, errors.New("record not found"))

	result, err := workoutUsecase.GetWorkoutByID(99, 2)

	assert.Error(t, err)
	assert.Equal(t, "data tidak ditemukan", err.Error())
	assert.Nil(t, result)
}

// ==========================================
// 3. UNIT TEST UNTUK FUNGSI CREATE
// ==========================================
func TestCreateWorkout_Success(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	val := validator.New()
	workoutUsecase := usecase.NewWorkoutUsecase(mockRepo, val)

	req := domain.WorkoutRequest{Name: "Push Day", Description: "Chest and Triceps"}

	// mock.Anything dipakai karena struct Workout di-generate di dalam usecase
	mockRepo.On("Create", mock.AnythingOfType("*domain.Workout")).Return(nil)

	result, err := workoutUsecase.CreateWorkout(2, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Push Day", result.Name)
	assert.Equal(t, 2, result.UserID)
}

// ==========================================
// 4. UNIT TEST UNTUK FUNGSI UPDATE (PUT)
// ==========================================
func TestUpdateWorkout_Success(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	val := validator.New()
	workoutUsecase := usecase.NewWorkoutUsecase(mockRepo, val)

	req := domain.WorkoutRequest{Name: "Pull Day Baru", Description: "Back and Biceps"}
	existingWorkout := &domain.Workout{ID: 1, UserID: 2, Name: "Pull Day Lama", Description: "Lama"}

	// Skenario: Cari data lama ketemu, lalu proses update berhasil
	mockRepo.On("GetByIDAndUserID", 1, 2).Return(existingWorkout, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Workout")).Return(nil)

	result, err := workoutUsecase.UpdateWorkout(1, 2, req)

	assert.NoError(t, err)
	assert.Equal(t, "Pull Day Baru", result.Name)
}

// ==========================================
// 5. UNIT TEST UNTUK FUNGSI DELETE
// ==========================================
func TestDeleteWorkout_Success(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	val := validator.New()
	workoutUsecase := usecase.NewWorkoutUsecase(mockRepo, val)

	existingWorkout := &domain.Workout{ID: 1, UserID: 2, Name: "Cardio"}

	// Skenario: Cari data ketemu, lalu hapus berhasil
	mockRepo.On("GetByIDAndUserID", 1, 2).Return(existingWorkout, nil)
	mockRepo.On("Delete", existingWorkout).Return(nil)

	result, err := workoutUsecase.DeleteWorkout(1, 2)

	assert.NoError(t, err)
	assert.Equal(t, "Cardio", result.Name)
}

func TestDeleteWorkout_NotFound(t *testing.T) {
	mockRepo := new(MockWorkoutRepository)
	val := validator.New()
	workoutUsecase := usecase.NewWorkoutUsecase(mockRepo, val)

	// Skenario: Cari data gagal, fungsi Delete tidak boleh dipanggil
	mockRepo.On("GetByIDAndUserID", 99, 2).Return(nil, errors.New("record not found"))

	result, err := workoutUsecase.DeleteWorkout(99, 2)

	assert.Error(t, err)
	assert.Equal(t, "data tidak ditemukan", err.Error())
	assert.Nil(t, result)
	mockRepo.AssertNotCalled(t, "Delete")
}
