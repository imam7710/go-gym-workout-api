package delivery

import (
	"gym-app/domain"
	"gym-app/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

// UserHandler adalah struct untuk mengikat dependency Usecase
type UserHandler struct {
	userUsecase domain.UserUsecase
}

// NewUserHandler berfungsi untuk mendaftarkan semua rute (routes) yang berkaitan dengan User
func NewUserHandler(e *echo.Echo, us domain.UserUsecase) {
	handler := &UserHandler{
		userUsecase: us,
	}

	// Mendaftarkan route API ke Echo framework
	api := e.Group("/api/users")
	api.POST("/register", handler.Register)
	api.POST("/login", handler.Login)
	api.GET("", handler.GetProfile, middleware.JWTMiddleware)
}

// Register adalah handler untuk endpoint POST /api/users/register
func (h *UserHandler) Register(c echo.Context) error {
	var req domain.UserRegisterRequest

	// 1. Bind JSON request ke struct
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "format request tidak valid",
		})
	}

	// 2. Panggil logika bisnis di Usecase
	user, err := h.userUsecase.Register(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	// 3. Kembalikan response sukses (status 201 Created)
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "berhasil mendaftarkan user baru",
		"data":    user,
	})
}

// Login adalah handler untuk endpoint POST /api/users/login
func (h *UserHandler) Login(c echo.Context) error {
	var req domain.UserLoginRequest

	// 1. Bind JSON request ke struct
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "format request tidak valid",
		})
	}

	// 2. Proses login ke usecase untuk mendapatkan JWT Token
	token, err := h.userUsecase.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
	}

	// 3. Kembalikan response sukses beserta token (status 200 OK)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "login berhasil",
		"token":   token,
	})
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	// Ambil user_id dari hasil ekstrak JWT Middleware tadi
	userID := c.Get("user_id").(int)

	user, bmi, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "berhasil mengambil profil user",
		"data":    user,
		"bmi":     bmi,
	})
}
