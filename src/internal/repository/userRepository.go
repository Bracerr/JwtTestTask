package repository

import (
	"JwtTestTask/src/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type UserRepositoryInterface interface {
	FindByGUID(guid string) (*domain.User, error)
	InsertUser(user domain.User) error
	UpdateUser(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	GetAll(page, limit int) ([]domain.User, int64, error)
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) FindByGUID(guid string) (*domain.User, error) {
	var user domain.User
	if err := repo.db.First(&user, "guid = ?", guid).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) InsertUser(user domain.User) error {
	return repo.db.Create(user).Error
}

func (repo *UserRepository) UpdateUser(user *domain.User) error {
	return repo.db.Save(user).Error
}

func (repo *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := repo.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetAll(page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	query := repo.db.Model(&domain.User{})

	query.Count(&total)
	query.Offset((page - 1) * limit).Limit(limit).Find(&users)

	return users, total, query.Error
}
