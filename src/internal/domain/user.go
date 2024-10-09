package domain

import "github.com/google/uuid"

type User struct {
	GUID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"guid"`
	RefreshToken *string   `json:"refresh_token" gorm:"type:text"`
	Email        string    `gorm:"unique" json:"email"`
}
