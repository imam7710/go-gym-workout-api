package domain

import "time"

type ExerciseLog struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ExerciseID int       `json:"exercise_id" gorm:"not null"`
	UserID     int       `json:"user_id" gorm:"not null"`
	SetCount   int       `json:"set_count" gorm:"not null"`
	RepCount   int       `json:"rep_count" gorm:"not null"`
	Weight     int       `json:"weight" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type ExerciseLogRequest struct {
	ExerciseID int `json:"exercise_id" validate:"required"`
	SetCount   int `json:"set_count" validate:"required,min=1"`
	RepCount   int `json:"rep_count" validate:"required,min=1"`
	Weight     int `json:"weight" validate:"required,min=0"`
}

type ExerciseLogRepository interface {
	Create(log *ExerciseLog) error
	GetByExerciseID(exerciseID int) ([]ExerciseLog, error)
}

type ExerciseLogUsecase interface {
	CreateLog(userID int, req ExerciseLogRequest) (*ExerciseLog, error)
	GetLogsByExerciseID(exerciseID int, userID int) ([]ExerciseLog, error)
}
