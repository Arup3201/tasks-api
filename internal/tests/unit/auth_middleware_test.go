package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Arup3201/gotasks/internal/middlewares"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret"
	issuer := "test-issuer"
	jwtSvc := utils.NewJWTService(secret, issuer)

	// Generate a valid token
	validToken, err := jwtSvc.GenerateToken("user-123", "test@example.com")
	require.NoError(t, err)

	// Generate an expired token
	expiredClaims := utils.Claims{
		UserID: "user-exp",
		Email:  "expired@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    issuer,
		},
	}
	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString([]byte(secret))
	require.NoError(t, err)

	// Generate a token with invalid signature
	invalidSigClaims := utils.Claims{
		UserID: "user-invalid",
		Email:  "invalid@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}
	invalidSigToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, invalidSigClaims).SignedString([]byte("wrong-secret"))
	require.NoError(t, err)

	cases := []struct {
		name              string
		authHeader        string
		wantStatus        int
		wantError         string
		wantUserID        string
		wantUserEmail     string
		expectNextHandler bool
	}{
		{
			name:              "valid token",
			authHeader:        "Bearer " + validToken,
			wantStatus:        http.StatusOK,
			wantUserID:        "user-123",
			wantUserEmail:     "test@example.com",
			expectNextHandler: true,
		},
		{
			name:              "missing authorization header",
			authHeader:        "",
			wantStatus:        http.StatusBadRequest,
			wantError:         "authorization header is empty",
			expectNextHandler: false,
		},
		{
			name:              "missing bearer prefix",
			authHeader:        validToken,
			wantStatus:        http.StatusBadRequest,
			wantError:         "invalid token format",
			expectNextHandler: false,
		},
		{
			name:              "wrong bearer prefix",
			authHeader:        "Basic " + validToken,
			wantStatus:        http.StatusBadRequest,
			wantError:         "invalid token format",
			expectNextHandler: false,
		},
		{
			name:              "bearer only",
			authHeader:        "Bearer",
			wantStatus:        http.StatusBadRequest,
			wantError:         "invalid token format",
			expectNextHandler: false,
		},
		{
			name:              "expired token",
			authHeader:        "Bearer " + expiredToken,
			wantStatus:        http.StatusBadRequest,
			wantError:         "token has expired",
			expectNextHandler: false,
		},
		{
			name:              "invalid signature",
			authHeader:        "Bearer " + invalidSigToken,
			wantStatus:        http.StatusBadRequest,
			wantError:         "invalid token",
			expectNextHandler: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nextHandlerCalled := false
			var capturedUserID string
			var capturedUserEmail string

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextHandlerCalled = true
				capturedUserID, _ = r.Context().Value("user_id").(string)
				capturedUserEmail, _ = r.Context().Value("user_email").(string)
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			rec := httptest.NewRecorder()

			// Create middleware with JWT service
			authMiddleware := middlewares.NewAuthMiddleware(jwtSvc)
			handler := authMiddleware.Required(nextHandler)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code, "unexpected status code")

			if tc.wantError != "" {
				assert.Contains(t, rec.Body.String(), tc.wantError, "expected error message not found")
			}

			assert.Equal(t, tc.expectNextHandler, nextHandlerCalled, "next handler call mismatch")

			if tc.expectNextHandler {
				assert.Equal(t, tc.wantUserID, capturedUserID, "user_id mismatch in context")
				assert.Equal(t, tc.wantUserEmail, capturedUserEmail, "user_email mismatch in context")
			}
		})
	}
}
