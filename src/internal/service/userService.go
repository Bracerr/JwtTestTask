package service

import (
	"JwtTestTask/src/internal/domain"
	"JwtTestTask/src/internal/payload/response"
	"JwtTestTask/src/internal/repository"
	"JwtTestTask/src/pkg/auth"
	"JwtTestTask/src/pkg/config"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/smtp"
	"time"
)

type UserService struct {
	repo         repository.UserRepositoryInterface
	tokenManager auth.JwtManagerInterface
}

type UserServiceInterface interface {
	SignIn(guid string, ip string) (response.JwtResponse, error)
	SignUp(email string) error
	RefreshTokens(accessToken string, refreshToken string, currentIp string) (response.JwtResponse, error)
	sendEmailWarning(email, oldIP, newIP string) error
	GetAll(page, limit int) ([]domain.User, int64, error)
}

func NewUserService(repo repository.UserRepositoryInterface, manager auth.JwtManagerInterface) *UserService {
	return &UserService{repo: repo, tokenManager: manager}
}

func (s *UserService) SignIn(guid string, ip string) (response.JwtResponse, error) {

	user, err := s.repo.FindByGUID(guid)
	if err != nil {
		return response.JwtResponse{}, fmt.Errorf("user not found")
	}

	accessToken, err := s.tokenManager.NewAccessToken(guid, ip)
	if err != nil {
		return response.JwtResponse{}, err
	}
	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return response.JwtResponse{}, err
	}

	if user.RefreshToken == nil {
		user.RefreshToken = new(string)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return response.JwtResponse{}, err
	}
	*user.RefreshToken = string(hash)

	if user.RefreshTokenExpiry == nil {
		user.RefreshTokenExpiry = new(time.Time)
	}
	expiryTime := time.Now().Add(s.tokenManager.GetRefreshDuration())
	*user.RefreshTokenExpiry = expiryTime

	err = s.repo.UpdateUser(user)
	tokens := response.JwtResponse{AccessToken: accessToken, RefreshToken: refreshToken}
	return tokens, nil
}

func (s *UserService) SignUp(email string) error {
	user := domain.User{
		GUID:               uuid.New(),
		RefreshToken:       nil,
		RefreshTokenExpiry: nil,
		Email:              email,
	}
	return s.repo.InsertUser(user)
}

func (s *UserService) RefreshTokens(accessToken string, refreshToken string, currentIp string) (response.JwtResponse, error) {
	claims, err := s.tokenManager.Parse(accessToken)
	if err != nil {
		return response.JwtResponse{}, err
	}

	user, err := s.repo.FindByGUID(claims.Subject)
	if err != nil {
		return response.JwtResponse{}, errors.New("user not found")
	}

	if user.RefreshTokenExpiry == nil || time.Now().After(*user.RefreshTokenExpiry) {
		return response.JwtResponse{}, fmt.Errorf("refresh token expired")
	}

	if user.RefreshToken == nil || bcrypt.CompareHashAndPassword([]byte(*user.RefreshToken), []byte(refreshToken)) != nil {
		return response.JwtResponse{}, errors.New("invalid refresh token")
	}

	if claims.IP != currentIp {
		err := s.sendEmailWarning(user.Email, claims.IP, currentIp)
		if err != nil {
			return response.JwtResponse{}, err
		}
		return response.JwtResponse{}, fmt.Errorf("invalid ip. Email warning")
	}

	newAccessToken, err := s.tokenManager.NewAccessToken(claims.Subject, claims.IP)
	if err != nil {
		return response.JwtResponse{}, err
	}

	newRefreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return response.JwtResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return response.JwtResponse{}, err
	}

	if user.RefreshToken == nil {
		user.RefreshToken = new(string)
	}
	*user.RefreshToken = string(hash)

	err = s.repo.UpdateUser(user)
	if err != nil {
		return response.JwtResponse{}, err
	}

	tokens := response.JwtResponse{AccessToken: newAccessToken, RefreshToken: newRefreshToken}
	return tokens, nil
}

func (s *UserService) sendEmailWarning(email, oldIP, newIP string) error {
	err := s.logout(email)
	if err != nil {
		return err
	}
	subject := "IP Address Change Warning"
	body := "Your IP address has changed from " + oldIP + " to " + newIP + "."
	message := []byte("Subject: " + subject + "\n\n" + body)

	smtpParams := config.GetSmtpParams()

	from := smtpParams.Username
	user := smtpParams.Username
	password := smtpParams.Password
	to := []string{email}

	addr := smtpParams.Host + ":" + smtpParams.Port
	host := smtpParams.Host

	plainAuth := smtp.PlainAuth("", user, password, host)
	err = smtp.SendMail(addr, plainAuth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) logout(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}
	*user.RefreshToken = ""
	err = s.repo.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetAll(page, limit int) ([]domain.User, int64, error) {
	return s.repo.GetAll(page, limit)
}
