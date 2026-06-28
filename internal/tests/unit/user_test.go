package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Arup3201/gotasks/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type mockUserStore struct {
	mock.Mock
}

func (m *mockUserStore) Create(ctx context.Context, id, name, email string, passwordHash []byte) error {
	args := m.Called(ctx, id, name, email, passwordHash)
	return args.Error(0)
}

func (m *mockUserStore) Get(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name           string
		email          string
		username       string
		password       string
		setupMock      func(store *mockUserStore)
		wantErr        error
		wantErrMessage string
		expectCreate   bool
		expectGet      bool
	}{
		{
			name:     "success",
			email:    "test@example.com",
			username: "Test User",
			password: "secret123",
			setupMock: func(store *mockUserStore) {
				store.
					On("Create", mock.Anything, mock.AnythingOfType("string"), "Test User", "test@example.com", mock.Anything).
					Return(nil)
				store.
					On("Get", mock.Anything, mock.AnythingOfType("string")).
					Return(&models.User{
						Email:     "test@example.com",
						Name:      "Test User",
						ID:        "generated-id",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
			},
			wantErr:      nil,
			expectCreate: true,
			expectGet:    true,
		},
		{
			name:      "invalid email",
			email:     "not-an-email",
			username:  "Test User",
			password:  "secret123",
			setupMock: func(store *mockUserStore) {},
			wantErr:   models.ErrInvalidEmail,
		},
		{
			name:      "empty name",
			email:     "test@example.com",
			username:  "   ",
			password:  "secret123",
			setupMock: func(store *mockUserStore) {},
			wantErr:   models.ErrEmptyName,
		},
		{
			name:     "store error",
			email:    "test@example.com",
			username: "Test User",
			password: "secret123",
			setupMock: func(store *mockUserStore) {
				store.
					On("Create", mock.Anything, mock.AnythingOfType("string"), "Test User", "test@example.com", mock.Anything).
					Return(errors.New("database failure"))
			},
			wantErrMessage: "database failure",
			expectCreate:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &mockUserStore{}
			if tc.setupMock != nil {
				tc.setupMock(store)
			}

			service := models.NewUserService(store)
			user, err := service.CreateUser(ctx, tc.email, tc.username, tc.password)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, user)
			} else if tc.wantErrMessage != "" {
				assert.ErrorContains(t, err, tc.wantErrMessage)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.email, user.Email)
				assert.Equal(t, tc.username, user.Name)
				assert.NotEmpty(t, user.ID)
				assert.False(t, user.CreatedAt.IsZero())
				assert.False(t, user.UpdatedAt.IsZero())
			}

			if tc.expectCreate {
				store.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("string"), tc.username, tc.email, mock.Anything)
			} else {
				store.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
			if tc.expectGet {
				store.AssertCalled(t, "Get", mock.Anything, mock.AnythingOfType("string"))
			} else {
				store.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
			}
			store.AssertExpectations(t)
		})
	}
}
func TestExchangeUserIDWithCredentials(t *testing.T) {
	ctx := context.Background()

	// hashed password for successful case
	hashed, err := bcrypt.GenerateFromPassword([]byte("secret123"), models.PasswordCost)
	if err != nil {
		t.Fatalf("bcrypt generate: %v", err)
	}

	cases := []struct {
		name      string
		email     string
		password  string
		setupMock func(store *mockUserStore)
		wantID    string
		wantErr   error
	}{
		{
			name:     "success",
			email:    "test@example.com",
			password: "secret123",
			setupMock: func(store *mockUserStore) {
				store.
					On("GetByEmail", mock.Anything, "test@example.com").
					Return(&models.User{
						ID:           "user-1",
						PasswordHash: hashed,
					}, nil)
			},
			wantID: "user-1",
		},
		{
			name:     "wrong password",
			email:    "test@example.com",
			password: "wrongpass",
			setupMock: func(store *mockUserStore) {
				store.
					On("GetByEmail", mock.Anything, "test@example.com").
					Return(&models.User{
						ID:           "user-1",
						PasswordHash: hashed,
					}, nil)
			},
			wantErr: models.ErrInvalidCredentials,
		},
		{
			name:     "store error",
			email:    "test@example.com",
			password: "secret123",
			setupMock: func(store *mockUserStore) {
				store.
					On("GetByEmail", mock.Anything, "test@example.com").
					Return((*models.User)(nil), errors.New("db failure"))
			},
			wantErr: models.ErrInvalidCredentials,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &mockUserStore{}
			if tc.setupMock != nil {
				tc.setupMock(store)
			}

			svc := models.NewUserService(store)
			user, err := svc.ExchangeUserWithCredentials(ctx, tc.email, tc.password)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantID, user.ID)
			}

			store.AssertExpectations(t)
		})
	}
}
