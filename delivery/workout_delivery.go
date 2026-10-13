package delivery

import (
	"gym-app/domain"
	"gym-app/middleware"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type WorkoutHandler struct {
	workoutUsecase domain.WorkoutUsecase
}

func NewWorkoutHandler(e *echo.Echo, workoutUsecase domain.WorkoutUsecase) {
	handler := &WorkoutHandler{
		workoutUsecase: workoutUsecase,
	}

	workoutGroup := e.Group("/api/workouts")
	workoutGroup.Use(middleware.JWTMiddleware)

	workoutGroup.POST("", handler.CreateWorkout)
	workoutGroup.GET("", handler.GetWorkouts)
	workoutGroup.GET("/:id", handler.GetWorkoutByID) // Tambahan: GET by ID
	workoutGroup.PUT("/:id", handler.UpdateWorkout)  // Tambahan: PUT
	workoutGroup.DELETE("/:id", handler.DeleteWorkout)
}

func extractUserID(c echo.Context) (int, bool) {
	rawID := c.Get("user_id")
	if rawID == nil {
		return 0, false
	}

	switch v := rawID.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case int64:
		return int(v), true
	default:
		return 0, false
	}
}

func (h *WorkoutHandler) CreateWorkout(c echo.Context) error {
	userID, ok := extractUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized: Gagal mendeteksi user ID dari token",
		})
	}

	var req domain.WorkoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Format request tidak valid",
		})
	}

	workout, err := h.workoutUsecase.CreateWorkout(userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Workout berhasil dibuat",
		"data":    workout,
	})
}

func (h *WorkoutHandler) GetWorkouts(c echo.Context) error {
	userID, ok := extractUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized: Gagal mendeteksi user ID dari token",
		})
	}

	workouts, err := h.workoutUsecase.GetWorkoutsByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal mengambil data workout",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Berhasil mengambil data workout",
		"data":    workouts,
	})
}

// Tambahan: Handler untuk GET Workout by ID
func (h *WorkoutHandler) GetWorkoutByID(c echo.Context) error {
	userID, ok := extractUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized: Gagal mendeteksi user ID dari token",
		})
	}

	workoutID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ID workout tidak valid",
		})
	}

	workout, err := h.workoutUsecase.GetWorkoutByID(workoutID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": err.Error(), // Mengembalikan 404 Not Found
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Berhasil mengambil data workout",
		"data":    workout,
	})
}

// Tambahan: Handler untuk UPDATE Workout
func (h *WorkoutHandler) UpdateWorkout(c echo.Context) error {
	userID, ok := extractUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized: Gagal mendeteksi user ID dari token",
		})
	}

	workoutID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ID workout tidak valid",
		})
	}

	var req domain.WorkoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Format request tidak valid",
		})
	}

	workout, err := h.workoutUsecase.UpdateWorkout(workoutID, userID, req)
	if err != nil {
		if err.Error() == "data tidak ditemukan" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Workout berhasil diperbarui",
		"data":    workout,
	})
}

// Perbaikan: Hapus workout membedakan 404 dan 500, serta return datanya
func (h *WorkoutHandler) DeleteWorkout(c echo.Context) error {
	userID, ok := extractUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized: Gagal mendeteksi user ID dari token",
		})
	}

	workoutID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ID workout tidak valid",
		})
	}

	// Tangkap data yang dihapus
	workout, err := h.workoutUsecase.DeleteWorkout(workoutID, userID)
	if err != nil {
		// Bedakan error 404 dengan 500
		if err.Error() == "data tidak ditemukan" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Workout berhasil dihapus",
		"data":    workout,
	})
}
