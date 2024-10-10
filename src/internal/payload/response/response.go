package response

import "JwtTestTask/src/internal/domain"

type ErrorResponse struct {
	Error string `json:"error"`
}

type JwtResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UsersResponse struct {
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
	Users []domain.User `json:"users"`
}
