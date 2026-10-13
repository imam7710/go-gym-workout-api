package domain

type Exercise struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	WorkoutID   int    `json:"workout_id" gorm:"not null"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
}

type ExerciseRequest struct {
	WorkoutID   int    `json:"workout_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// Tambahkan dua interface ini di bawahnya:
type ExerciseRepository interface {
	Create(exercise *Exercise) error
	GetByWorkoutID(workoutID int) ([]Exercise, error)
	GetByID(id int) (*Exercise, error)
	Delete(exercise *Exercise) error
}

type ExerciseUsecase interface {
	CreateExercise(userID int, req ExerciseRequest) (*Exercise, error)
	GetExercisesByWorkoutID(workoutID int) ([]Exercise, error)
	DeleteExercise(id int, userID int) (*Exercise, error)
}
