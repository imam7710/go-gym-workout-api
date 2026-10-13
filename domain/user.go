package domain

// User Struct merepresentasikan tabel users di database
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"-"` // Tanda "-" memastikan password tidak ikut ter-render di response JSON
	Weight   int    `json:"weight"`
	Height   int    `json:"height"`
}

// UserRegisterRequest adalah DTO (Data Transfer Object) untuk request register
type UserRegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Weight   int    `json:"weight" validate:"required"`
	Height   int    `json:"height" validate:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserRepository adalah kontrak untuk akses ke database (PostgreSQL)
type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByID(id int) (*User, error)
}

// UserUsecase adalah kontrak untuk business logic
type UserUsecase interface {
	Register(req UserRegisterRequest) (*User, error)
	Login(email, password string) (string, error) // Return JWT Token
	GetProfile(id int) (*User, float64, error)    // Return User dan BMI (dari 3rd party API)
}
