package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type DBUser struct {
	ID                   string `gorm:"primaryKey"`
	Name                 string
	Email                string `gorm:"unique"`
	IsCompleted          bool
	PasswordHash         []byte
	CreatedAt, UpdatedAt time.Time
}

type UserStore struct {
	db *gorm.DB
}

func (us *UserStore) Create(ctx context.Context,
	id, name, email string,
	passwordHash []byte) error {

	user := DBUser{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	err := gorm.G[DBUser](us.db).Create(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserStore) Get(ctx context.Context,
	id string) (*User, error) {

	userRow, err := gorm.G[DBUser](us.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           userRow.ID,
		Name:         userRow.Name,
		Email:        userRow.Email,
		IsCompleted:  userRow.IsCompleted,
		PasswordHash: userRow.PasswordHash,
		CreatedAt:    userRow.CreatedAt,
		UpdatedAt:    userRow.UpdatedAt,
	}, nil
}
