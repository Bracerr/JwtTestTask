package main

import (
	_ "JwtTestTask/docs"
	"JwtTestTask/src/internal/domain"
	"JwtTestTask/src/internal/repository"
	"JwtTestTask/src/internal/routing"
	"JwtTestTask/src/internal/service"
	"JwtTestTask/src/pkg/auth"
	"JwtTestTask/src/pkg/config"
	"JwtTestTask/src/pkg/database"
	"JwtTestTask/src/pkg/logger"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	logger.Init()
	config.Init()

	dbModel := config.GetDbParams()
	db := database.NewClient(dbModel)
	logger.Log.Infoln("Database connection established")

	jwtModel := config.GetJwtParams()
	jwtManager, err := auth.NewManager(jwtModel.SigningKey, jwtModel.AccessDuration, jwtModel.RefreshDuration)

	if err != nil {
		logger.Log.Errorln(err.Error())
	}

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, jwtManager)

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		logger.Log.Fatal("Ошибка миграции:", err)
	} else {
		logger.Log.Infoln("Успешная миграция.")
	}

	e := echo.New()
	routing.SetupUserRoute(e, userService)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(config.GetServerParams().ServerHost))
}
