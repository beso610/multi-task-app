package repository

import (
	"encoding/base64"
	"errors"

	"github.com/google/uuid"
	"github.com/product/multi-task-app/server/utils"
	"gorm.io/gorm"
)

// CreateUserArgs ユーザー作成時に必要な情報
type CreateUserArgs struct {
	MailAddress string
	Password    string
	UserName    string
}

// UpdateUserArgs ユーザー更新時に必要な情報
type UpdateUserArgs struct {
	MailAddress string
	Password    string
	UserName    string
}

// UserRepository ユーザーリポジトリ
type UserRepository interface {
	CreateUser(args CreateUserArgs) (*User, error)
	UpdateUser(uid uuid.UUID, args UpdateUserArgs) (*User, error)
	GetUser(id uuid.UUID) (*User, error)
	GetUsers() ([]*User, error)
	GetUserByMailAddress(mailAddress string) (*User, error)
}

func (repo *DBRepository) CreateUser(args CreateUserArgs) (*User, error) {
	uid := uuid.Must(uuid.NewUUID())

	salt := utils.NewSalt64()
	password := utils.HashPassword(args.Password, salt)

	user := &User{
		ID:          uid,
		Password:    base64.RawURLEncoding.EncodeToString(password),
		Salt:        base64.RawURLEncoding.EncodeToString(salt),
		MailAddress: args.MailAddress,
		UserName:    args.UserName,
	}

	err := repo.DB.Create(user).Error

	return user, err
}

func (repo *DBRepository) UpdateUser(uid uuid.UUID, args UpdateUserArgs) (*User, error) {
	updateData := User{}
	if args.MailAddress != "" {
		updateData.MailAddress = args.MailAddress
	}
	if args.Password != "" {
		salt := utils.NewSalt64()
		password := utils.HashPassword(args.Password, salt)
		updateData.Password = base64.RawURLEncoding.EncodeToString(password)
		updateData.Salt = base64.RawURLEncoding.EncodeToString(salt)
	}
	if args.UserName != "" {
		updateData.UserName = args.UserName
	}
	err := repo.DB.Model(&User{}).Where("id = ?", uid).Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	return repo.GetUser(uid)
}

func (repo *DBRepository) GetUser(id uuid.UUID) (*User, error) {
	user := &User{}
	err := repo.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (repo *DBRepository) GetUsers() ([]*User, error) {
	users := make([]*User, 0)
	if err := repo.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *DBRepository) GetUserByMailAddress(mailAddress string) (*User, error) {
	user := &User{}
	err := repo.DB.Where(&User{MailAddress: mailAddress}).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
