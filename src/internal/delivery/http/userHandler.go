package http

import (
	"JwtTestTask/src/internal/domain"
	"JwtTestTask/src/internal/payload/response"
	"JwtTestTask/src/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

type UserHandlerInterface interface {
	UserSignIn(c echo.Context) error
	UserSignUp(c echo.Context) error
	RefreshTokens(c echo.Context) error
	GetAll(c echo.Context) error
}

// UserSignIn godoc
// @Summary User Sign In
// @Description Выдача access & refresh токенов по GUID user
// @Tags users
// @Accept json
// @Produce json
// @Param guid query string true "User GUID"
// @Success 200 {object} response.JwtResponse "Successful response"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Router /signIn [post]
func (h *UserHandler) UserSignIn(c echo.Context) error {
	guid := c.QueryParam("guid")
	ip := c.Request().RemoteAddr
	tokens, err := h.service.SignIn(guid, ip)
	if err != nil {
		errorResponse := response.ErrorResponse{Error: err.Error()}
		return c.JSON(http.StatusNotFound, errorResponse)
	}

	return c.JSON(http.StatusOK, tokens)
}

// UserSignUp godoc
// @Summary User Sign Up
// @Description Создание пользователя по email
// @Tags users
// @Accept json
// @Produce json
// @Param email query string true "User Email"
// @Success 201 {object} nil "User created successfully"
// @Failure 400 {object} response.ErrorResponse "Email is required"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Router /signUp [post]
func (h *UserHandler) UserSignUp(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		errorResponse := response.ErrorResponse{Error: "email is required"}
		return c.JSON(http.StatusBadRequest, errorResponse)
	}
	err := h.service.SignUp(email)
	if err != nil {
		errorResponse := response.ErrorResponse{Error: err.Error()}
		return c.JSON(http.StatusNotFound, errorResponse)
	}
	return c.NoContent(http.StatusCreated)
}

// RefreshTokens godoc
// @Summary Refresh JWT Tokens
// @Description Обновление токенов по паре access & refresh tokens.
// @Description При смене ip высылается email warning на почту указанную при создании и refresh токен в базе обнуляется.
// @Description Токены были перенесены из headers в body для удобства отладки и проверки задания
// @Tags users
// @Accept json
// @Produce json
// @Param tokensRequest body response.JwtResponse true "Tokens Request"
// @Success 200 {object} response.JwtResponse "Successful response with new tokens"
// @Failure 400 {object} response.ErrorResponse "Invalid request or tokens"
// @Router /refresh [post]
func (h *UserHandler) RefreshTokens(c echo.Context) error {
	var tokensRequest response.JwtResponse
	ip := c.Request().RemoteAddr

	if err := c.Bind(&tokensRequest); err != nil {
		errorResponse := response.ErrorResponse{Error: err.Error()}
		return c.JSON(http.StatusBadRequest, errorResponse)
	}
	tokens, err := h.service.RefreshTokens(tokensRequest.AccessToken, tokensRequest.RefreshToken, ip)
	if err != nil {
		errorResponse := response.ErrorResponse{Error: err.Error()}
		return c.JSON(http.StatusBadRequest, errorResponse)
	}
	return c.JSON(http.StatusOK, tokens)
}

// GetAll godoc
// @Summary Get All Users
// @Description Получение всех пользователей с пагинацией. Вспомогательный эндпоинт для более удобного тестирования
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of users per page" default(10)
// @Success 200 {object} response.UsersResponse "Successful response with user list"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /getAll [get]
func (h *UserHandler) GetAll(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	users, total, err := h.service.GetAll(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
	}

	var userResponse []domain.User

	for _, user := range users {
		userResponse = append(userResponse, domain.User{
			GUID:         user.GUID,
			RefreshToken: user.RefreshToken,
			Email:        user.Email,
		})
	}

	usersResponse := response.UsersResponse{
		Total: int(total),
		Page:  page,
		Limit: limit,
		Users: userResponse,
	}

	return c.JSON(http.StatusOK, usersResponse)
}
