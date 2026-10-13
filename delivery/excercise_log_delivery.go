package delivery

import (
	"gym-app/domain"
	"gym-app/middleware"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ExerciseLogHandler struct {
	logUsecase domain.ExerciseLogUsecase
}

func NewExerciseLogHandler(e *echo.Echo, logUsecase domain.ExerciseLogUsecase) {
	handler := &ExerciseLogHandler{
		logUsecase: logUsecase,
	}

	logGroup := e.Group("/api/logs")
	logGroup.Use(middleware.JWTMiddleware)

	logGroup.POST("", handler.CreateLog)
	logGroup.GET("/exercise/:exercise_id", handler.GetLogs)
}

// CreateLog godoc
// @Summary Catat Log Latihan (Reps & Weight)
// @Description Mencatat repetisi dan beban untuk exercise tertentu
// @Tags Logs
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param request body domain.ExerciseLogRequest true "Data Log"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /logs [post]
func (h *ExerciseLogHandler) CreateLog(c echo.Context) error {
	userID, ok := extractUserID(c) // Menggunakan fungsi extractUserID dari workout_delivery.go
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized: Gagal mendeteksi user ID dari token",
		})
	}

	var req domain.ExerciseLogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Format request tidak valid",
		})
	}

	log, err := h.logUsecase.CreateLog(userID, req)
	if err != nil {
		if err.Error() == "unauthorized: anda tidak memiliki akses ke exercise ini" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Log berhasil dicatat",
		"data":    log,
	})
}

// GetLogs godoc
// @Summary Lihat Histori Log
// @Description Melihat riwayat repetisi dan beban pada exercise tertentu
// @Tags Logs
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param exercise_id path int true "ID Exercise"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /logs/exercise/{exercise_id} [get]
func (h *ExerciseLogHandler) GetLogs(c echo.Context) error {
	userID, ok := extractUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized: Gagal mendeteksi user ID dari token",
		})
	}

	exerciseID, err := strconv.Atoi(c.Param("exercise_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ID exercise tidak valid",
		})
	}

	logs, err := h.logUsecase.GetLogsByExerciseID(exerciseID, userID)
	if err != nil {
		if err.Error() == "unauthorized: anda tidak memiliki akses ke exercise ini" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Berhasil mengambil histori log",
		"data":    logs,
	})
}
