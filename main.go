package main

import (
	"gym-app/config"
	"gym-app/delivery"
	_ "gym-app/docs"
	"gym-app/repository"
	"gym-app/usecase"
	"log"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// @title Gym App API
// @version 1.0
// @description API Backend untuk aplikasi Gym berbasis Clean Architecture menggunakan Go, Echo, dan GORM.
// @host localhost:8080
// @BasePath /api
func main() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic("Gagal membaca file .env: " + err.Error())
	}
	db := config.InitDB()
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	val := validator.New()

	// Injeksi dependensi
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, val)
	delivery.NewUserHandler(e, userUsecase)

	workoutRepo := repository.NewWorkoutRepository(db)
	workoutUsecase := usecase.NewWorkoutUsecase(workoutRepo, val)
	delivery.NewWorkoutHandler(e, workoutUsecase)

	exerciseRepo := repository.NewExerciseRepository(db)
	exerciseUsecase := usecase.NewExerciseUsecase(exerciseRepo, workoutRepo, val)
	delivery.NewExerciseHandler(e, exerciseUsecase)

	logRepo := repository.NewExerciseLogRepository(db)
	logUsecase := usecase.NewExerciseLogUsecase(logRepo, exerciseRepo, workoutRepo, val)
	delivery.NewExerciseLogHandler(e, logUsecase)

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
