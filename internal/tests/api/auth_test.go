package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	controller "github.com/Arup3201/gotasks/internal/controllers"
	"github.com/Arup3201/gotasks/internal/models"
	"github.com/Arup3201/gotasks/internal/testutils"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type registerResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type loginResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	User        struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
}

type authTestEnv struct {
	controller *controller.AuthController
	userSvc    *models.UserService
	jwtSvc     *utils.JWTService
}

func setupAuthTestEnv(t *testing.T) *authTestEnv {
	t.Helper()

	ctx := context.Background()
	pg, err := testutils.CreatePostgresContainer(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := pg.Terminate(ctx)
		require.NoError(t, err)
	})

	db, err := gorm.Open(postgres.Open(pg.ConnectionString), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.DBUser{})
	require.NoError(t, err)

	userSvc := models.NewUserService(models.NewUserStore(db))
	jwtSvc := utils.NewJWTService("test-secret", "test-issuer")

	return &authTestEnv{
		controller: controller.NewAuthController(userSvc, jwtSvc),
		userSvc:    userSvc,
		jwtSvc:     jwtSvc,
	}
}

func TestRegisterEndpoint(t *testing.T) {
	env := setupAuthTestEnv(t)

	cases := []struct {
		name       string
		body       string
		wantStatus int
		wantEmail  string
		wantName   string
		wantError  string
	}{
		{
			name:       "success",
			body:       `{"name":"Alice Example","email":"alice@example.com","password":"secret123"}`,
			wantStatus: http.StatusCreated,
			wantEmail:  "alice@example.com",
			wantName:   "Alice Example",
		},
		{
			name:       "invalid email",
			body:       `{"name":"Alice Example","email":"bad-email","password":"secret123"}`,
			wantStatus: http.StatusInternalServerError,
			wantError:  "email address is invalid",
		},
		{
			name:       "empty name",
			body:       `{"name":"   ","email":"alice@example.com","password":"secret123"}`,
			wantStatus: http.StatusInternalServerError,
			wantError:  "name cannot be empty",
		},
		{
			name:       "malformed json",
			body:       `{"name":"Alice","email":"alice@example.com","password":"secret123"`,
			wantStatus: http.StatusBadRequest,
			wantError:  "json parse error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tc.body))
			rec := httptest.NewRecorder()

			env.controller.Register(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)

			if tc.wantError != "" {
				assert.Contains(t, rec.Body.String(), tc.wantError)
				return
			}

			var resp registerResponse
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			assert.Equal(t, tc.wantEmail, resp.Email)
			assert.Equal(t, tc.wantName, resp.Name)
			assert.NotEmpty(t, resp.ID)
			assert.False(t, resp.CreatedAt.IsZero())
		})
	}
}

func TestLoginEndpoint(t *testing.T) {
	env := setupAuthTestEnv(t)

	// Prepare a registered user for login tests.
	user, err := env.userSvc.CreateUser(context.Background(), "bob@example.com", "Bob Example", "secret123")
	require.NoError(t, err)

	cases := []struct {
		name            string
		setup           func(t *testing.T)
		body            string
		wantStatus      int
		wantError       string
		wantUserID      string
		wantUserEmail   string
		expectJWTClaims bool
	}{
		{
			name:            "success",
			body:            `{"email":"bob@example.com","password":"secret123"}`,
			wantStatus:      http.StatusOK,
			wantUserID:      user.ID,
			wantUserEmail:   user.Email,
			expectJWTClaims: true,
		},
		{
			name:       "wrong password",
			body:       `{"email":"bob@example.com","password":"badpass"}`,
			wantStatus: http.StatusInternalServerError,
			wantError:  "invalid login credentials",
		},
		{
			name:       "unknown email",
			body:       `{"email":"missing@example.com","password":"secret123"}`,
			wantStatus: http.StatusInternalServerError,
			wantError:  "invalid login credentials",
		},
		{
			name:       "malformed json",
			body:       `{"email":"bob@example.com","password":"secret123"`,
			wantStatus: http.StatusBadRequest,
			wantError:  "json parse error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tc.body))
			rec := httptest.NewRecorder()

			env.controller.Login(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)

			if tc.wantError != "" {
				assert.Contains(t, rec.Body.String(), tc.wantError)
				return
			}

			var resp loginResponse
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			assert.NotEmpty(t, resp.AccessToken)
			assert.Equal(t, tc.wantUserID, resp.User.ID)
			assert.Equal(t, tc.wantUserEmail, resp.User.Email)
			assert.Equal(t, "Bob Example", resp.User.Name)
			assert.WithinDuration(t, time.Now().Add(utils.TOKEN_EXPIRES_IN), resp.ExpiresAt, time.Minute)

			if tc.expectJWTClaims {
				claims, err := env.jwtSvc.ValidateToken(resp.AccessToken)
				require.NoError(t, err)
				assert.Equal(t, tc.wantUserID, claims.UserID)
				assert.Equal(t, tc.wantUserEmail, claims.Email)
				assert.Equal(t, "test-issuer", claims.Issuer)
				assert.WithinDuration(t, time.Now().Add(utils.TOKEN_EXPIRES_IN), claims.ExpiresAt.Time, time.Minute)
			}
		})
	}
}
