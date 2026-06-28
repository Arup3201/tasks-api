package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Arup3201/gotasks/internal/controllers"
	"github.com/Arup3201/gotasks/internal/models"
	"github.com/Arup3201/gotasks/internal/testutils"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

type AuthTestSuite struct {
	suite.Suite
	controller *controllers.AuthController
	userSvc    *models.UserService
	jwtSvc     *utils.JWTService
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) SetupSuite() {

	ctx := context.Background()
	pg, err := testutils.CreatePostgresContainer(ctx)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		s.Require().NoError(pg.Terminate(ctx), "could not terminate postgres container")
	})

	db, err := gorm.Open(postgres.Open(pg.ConnectionString), &gorm.Config{})
	s.Require().NoError(err)

	s.Require().NoError(db.AutoMigrate(&models.User{}))

	s.userSvc = models.NewUserService(models.NewUserStore(db))
	s.jwtSvc = utils.NewJWTService("test-secret", "test-issuer")

	s.controller = controllers.NewAuthController(s.userSvc, s.jwtSvc)
}

func (s *AuthTestSuite) TestRegisterEndpoint() {
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
		s.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tc.body))
			rec := httptest.NewRecorder()

			s.controller.Register(rec, req)

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

func (s *AuthTestSuite) TestLoginEndpoint() {
	// Prepare a registered user for login tests.
	user, err := s.userSvc.CreateUser(context.Background(), "bob@example.com", "Bob Example", "secret123")
	require.NoError(s.T(), err)

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
		s.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tc.body))
			rec := httptest.NewRecorder()

			s.controller.Login(rec, req)

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
				claims, err := s.jwtSvc.ValidateToken(resp.AccessToken)
				require.NoError(t, err)
				assert.Equal(t, tc.wantUserID, claims.UserID)
				assert.Equal(t, tc.wantUserEmail, claims.Email)
				assert.Equal(t, "test-issuer", claims.Issuer)
				assert.WithinDuration(t, time.Now().Add(utils.TOKEN_EXPIRES_IN), claims.ExpiresAt.Time, time.Minute)
			}
		})
	}
}
