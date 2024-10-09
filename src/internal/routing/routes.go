package routing

import (
	"JwtTestTask/src/internal/delivery/http"
	"JwtTestTask/src/internal/service"
	"github.com/labstack/echo/v4"
)

func SetupUserRoute(e *echo.Echo, userService *service.UserService) {
	userHandler := http.NewUserHandler(userService)

	e.POST("/signIn", userHandler.UserSignIn)
	e.POST("/signUp", userHandler.UserSignUp)
	e.POST("/refresh", userHandler.RefreshTokens)
	e.GET("/getAll", userHandler.GetAll)
}
