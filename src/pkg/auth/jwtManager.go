package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type JwtManager struct {
	signingKey string
	duration   time.Duration
}

type CustomClaims struct {
	IP string `json:"ip"`
	jwt.StandardClaims
}

type JwtManagerInterface interface {
	NewAccessToken(guid string, ip string) (string, error)
	NewRefreshToken() (string, error)
	Parse(accessToken string) (*CustomClaims, error)
}

func NewManager(signingKey string, jwtDuration time.Duration) (*JwtManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}
	if jwtDuration <= 0 {
		return nil, errors.New("invalid duration")
	}
	return &JwtManager{signingKey: signingKey, duration: jwtDuration}, nil
}

func (m *JwtManager) NewAccessToken(guid string, ip string) (string, error) {
	claims := CustomClaims{
		IP: ip,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.duration).Unix(),
			Subject:   guid,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(m.signingKey))
}

func (m *JwtManager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(b)
	return refreshToken, nil
}

func (m *JwtManager) Parse(accessToken string) (*CustomClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error getting user claims from token: %v", token.Claims)
	}

	customClaims := &CustomClaims{
		IP: claims["ip"].(string),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(claims["exp"].(float64)),
			Subject:   claims["sub"].(string),
		},
	}

	return customClaims, nil
}
