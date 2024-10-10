package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	GUID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"guid"`
	RefreshToken       *string    `json:"refresh_token" gorm:"type:text"`
	RefreshTokenExpiry *time.Time `json:"refresh_token_expiry" gorm:"type:timestamp"`
	Email              string     `gorm:"unique" json:"email"`
}
