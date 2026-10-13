package delivery

import (
	"gym-app/domain"
	"gym-app/middleware"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ExerciseHandler struct {
	exerciseUsecase domain.ExerciseUsecase
}

func NewExerciseHandler(e *echo.Echo, exerciseUsecase domain.ExerciseUsecase) {
	handler := &ExerciseHandler{
		exerciseUsecase: exerciseUsecase,
	}

	exerciseGroup := e.Group("/api/exercises")
	exerciseGroup.Use(middleware.JWTMiddleware)

	exerciseGroup.POST("", handler.CreateExercise)
	exerciseGroup.GET("/workout/:workout_id", handler.GetExercises)
	exerciseGroup.DELETE("/:id", handler.DeleteExercise)
}

// CreateExercise godoc
// @Summary Tambah Exercise baru
// @Description Menambahkan latihan baru pada workout (hanya jika workout milik user yang login)
// @Tags Exercises
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param request body domain.ExerciseRequest true "Data Exercise"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /exercises [post]
func (h *ExerciseHandler) CreateExercise(c echo.Context) error {
	// Ambil User ID untuk memverifikasi kepemilikan Workout
	userID, ok := extractUserID(c) // Berbagi fungsi dengan workout_delivery.go
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized: Gagal mendeteksi user ID dari token",
		})
	}

	var req domain.ExerciseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Format request tidak valid",
		})
	}

	// Teruskan User ID ke usecase
	exercise, err := h.exerciseUsecase.CreateExercise(userID, req)
	if err != nil {
		if err.Error() == "unauthorized: workout tidak ditemukan atau bukan milik anda" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Exercise berhasil ditambahkan",
		"data":    exercise,
	})
}

// GetExercises godoc
// @Summary Lihat Daftar Exercise
// @Description Melihat daftar exercise berdasarkan ID workout
// @Tags Exercises
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param workout_id path int true "ID Workout"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /exercises/workout/{workout_id} [get]
func (h *ExerciseHandler) GetExercises(c echo.Context) error {
	workoutID, err := strconv.Atoi(c.Param("workout_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ID workout tidak valid",
		})
	}

	exercises, err := h.exerciseUsecase.GetExercisesByWorkoutID(workoutID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal mengambil data exercise",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Berhasil mengambil data exercise",
		"data":    exercises,
	})
}

// DeleteExercise godoc
// @Summary Hapus Exercise
// @Description Menghapus data exercise berdasarkan ID (hanya jika workout milik user yang login)
// @Tags Exercises
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "ID Exercise"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /exercises/{id} [delete]
func (h *ExerciseHandler) DeleteExercise(c echo.Context) error {
	userID, ok := extractUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized: Gagal mendeteksi user ID dari token",
		})
	}

	exerciseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ID exercise tidak valid",
		})
	}

	// Teruskan ID dan User ID ke Usecase
	exercise, err := h.exerciseUsecase.DeleteExercise(exerciseID, userID)
	if err != nil {
		if err.Error() == "data exercise tidak ditemukan" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": err.Error(), // 404
			})
		}
		if err.Error() == "unauthorized: tidak dapat menghapus exercise pada workout orang lain" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": err.Error(), // 403
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(), // 500
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Exercise berhasil dihapus",
		"data":    exercise,
	})
}
