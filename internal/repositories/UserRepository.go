package repositories

import (
	"avito-tech-go/internal/domain"
	"errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(*domain.User) error
	UpdateUser(*domain.User) error
	GetUserByID(uint) (*domain.User, error)
	DeleteUser(*domain.User) error
	GetUserByName(username string) (*domain.User, error)
	ExistsByUsername(username string) (bool, error)
	ChangeCoins(userID int, delta int) error
	GetUsernamesByIDs(ids []uint) (map[uint]string, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) CreateUser(user *domain.User) error {
	return u.db.Create(user).Error
}

func (u *userRepository) UpdateUser(user *domain.User) error {
	return u.db.Save(user).Error
}

func (u *userRepository) GetUserByID(ID uint) (*domain.User, error) {
	var user domain.User
	err := u.db.Where("ID = ?", ID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (u *userRepository) GetUserByName(username string) (*domain.User, error) {
	var user domain.User
	err := u.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (u *userRepository) DeleteUser(user *domain.User) error {
	return u.db.Delete(user).Error
}

func (u *userRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := u.db.Model(&domain.User{}).
		Where("username = ?", username).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, err
}

func (u *userRepository) ChangeCoins(userID int, delta int) error {
	res := u.db.Model(&domain.User{}).
		Where("ID = ?", userID).
		Update("coins", gorm.Expr("coins + ?", delta))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("no rows affected (users not found?)")
	}
	return nil
}

func (u *userRepository) GetUsernamesByIDs(ids []uint) (map[uint]string, error) {
	var users []domain.User
	if err := u.db.Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	result := make(map[uint]string, len(users))
	for _, user := range users {
		result[user.ID] = user.Username
	}
	return result, nil
}
