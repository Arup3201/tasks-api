package models

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex, _ = regexp.Compile(
	`^[a-zA-Z0-9]+([._-][0-9a-zA-Z]+)*@[a-zA-Z0-9]+([.-][0-9a-zA-Z]+)*\.[a-zA-Z]{2,}$`,
)

const (
	BCRYPT_PASSWORD_COST_DEFAULT = 14
)

var (
	PasswordCost = BCRYPT_PASSWORD_COST_DEFAULT
)

var (
	ErrInvalidEmail       = errors.New("email address is invalid")
	ErrEmptyName          = errors.New("name cannot be empty")
	ErrInvalidCredentials = errors.New("invalid login credentials")
)

type User struct {
	ID, Name, Email      string
	IsCompleted          bool
	PasswordHash         []byte
	CreatedAt, UpdatedAt time.Time
}

type UserStoreInterface interface {
	Create(ctx context.Context,
		id, name, email string,
		passwordHash []byte) error
	Get(ctx context.Context,
		id string) (*User, error)
	GetByEmail(ctx context.Context,
		email string) (*User, error)
}

type UserService struct {
	store UserStoreInterface
}

func NewUserService(store UserStoreInterface) *UserService {
	return &UserService{store: store}
}

func (us *UserService) CreateUser(ctx context.Context,
	email, name, password string) (*User, error) {
	var err error

	if match := emailRegex.Find([]byte(email)); match == nil {
		return nil, ErrInvalidEmail
	}

	if strings.Trim(name, " ") == "" {
		return nil, ErrEmptyName
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),
		PasswordCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt generate from password: %w", err)
	}

	id := uuid.NewString()
	err = us.store.Create(ctx, id, name, email, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("store create: %w", err)
	}

	var user *User
	user, err = us.store.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("store get: %w", err)
	}

	return user, nil
}

func (us *UserService) ExchangeUserIDWithCredentials(ctx context.Context,
	email, password string) (string, error) {

	user, err := us.store.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return user.ID, nil
}
