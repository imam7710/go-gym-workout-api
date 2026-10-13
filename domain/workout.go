package domain

type Workout struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
	UserID      int    `json:"user_id" gorm:"not null"`
	// Tambahan: Relasi ke tabel exercises agar terbawa saat GET Workout
	Exercises []Exercise `json:"exercises" gorm:"foreignKey:WorkoutID"`
}

type WorkoutRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type WorkoutRepository interface {
	Create(workout *Workout) error
	GetByUserID(userID int) ([]Workout, error)
	GetByIDAndUserID(id int, userID int) (*Workout, error)
	Update(workout *Workout) error
	Delete(workout *Workout) error
}

type WorkoutUsecase interface {
	CreateWorkout(userID int, req WorkoutRequest) (*Workout, error)
	GetWorkoutsByUserID(userID int) ([]Workout, error)
	GetWorkoutByID(id int, userID int) (*Workout, error)
	UpdateWorkout(id int, userID int, req WorkoutRequest) (*Workout, error)
	DeleteWorkout(id int, userID int) (*Workout, error)
}
