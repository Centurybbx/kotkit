package repository

import (
	"errors"
	"github.com/Tiktok-Lite/kotkit/internal/db"
	"github.com/Tiktok-Lite/kotkit/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	UpdateByUsername(username string, updatedUser *model.User) error
	QueryUserByID(id int64) (*model.User, error)
	QueryUserByName(name string) (*model.User, error)
	QueryUserByRelation(userID, followerID int64) (bool, error)
}

type userRepository struct {
	*Repository
}

func NewUserRepository(r *Repository) UserRepository {
	return &userRepository{
		Repository: r,
	}
}

func (r *userRepository) Create(user *model.User) error {
	// TODO(century): add error info to log
	if err := r.db.Create(user).Error; err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (r *userRepository) Update(user *model.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return errors.New("failed to update user")
	}

	return nil
}

func (r *userRepository) UpdateByUsername(username string, updatedUser *model.User) error {
	// 构建更新条件
	condition := map[string]interface{}{
		"Name": username,
	}
	// 执行更新操作
	if err := r.db.Model(&model.User{}).Where(condition).Updates(updatedUser).Error; err != nil {
		return errors.New("failed to update user by username")
	}

	return nil
}

func (r *userRepository) QueryUserByID(id int64) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) QueryUserByName(name string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("name = ?", name).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New("failed to query user by name")
	}
	return &user, nil
}

func (r *userRepository) QueryUserByRelation(userID, followerID int64) (bool, error) {
	var count int64
	err := db.DB().Raw("SELECT COUNT(*) FROM user_relations WHERE user_id = ? AND follower_id = ?", userID, followerID).Count(&count).Error
	if err != nil {
		return false, errors.New("failed to query user by relation")
	}
	return count > 0, nil
}
