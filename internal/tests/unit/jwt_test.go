package unit

import (
	"os"
	"testing"
	"time"

	"github.com/Arup3201/gotask/internal/config"
	"github.com/Arup3201/gotask/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTValidate(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_ISSUER", "test-issuer")
	config := config.Load()

	svc := utils.NewJWTService(config)

	validToken, err := svc.GenerateToken("user-1", "user@example.com")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	// expired token: create manually with past expiry
	expiredClaims := utils.Claims{
		UserID: "user-exp",
		Email:  "exp@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    config.JWT.Issuer,
		},
	}
	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString([]byte(config.JWT.Secret))
	if err != nil {
		t.Fatalf("create expired token: %v", err)
	}

	// invalid signature token
	invalidToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.Claims{
		UserID: "user-x",
		Email:  "x@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.JWT.Issuer,
		},
	}).SignedString([]byte("other-secret"))
	if err != nil {
		t.Fatalf("create invalid token: %v", err)
	}

	cases := []struct {
		name    string
		svc     *utils.JWTService
		token   string
		wantErr error
		wantID  string
	}{
		{name: "valid", svc: svc, token: validToken, wantErr: nil, wantID: "user-1"},
		{name: "expired", svc: svc, token: expiredToken, wantErr: utils.ErrExpiredToken},
		{name: "invalid signature", svc: svc, token: invalidToken, wantErr: utils.ErrInvalidToken},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := tc.svc.ValidateToken(tc.token)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantID, claims.UserID)
			}
		})
	}
}

func TestJWTRefresh(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_ISSUER", "test-issuer")
	config := config.Load()

	svc := utils.NewJWTService(config)

	token, err := svc.GenerateToken("user-refresh", "ref@example.com")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	newToken, err := svc.RefreshToken(token)
	if err != nil {
		t.Fatalf("refresh token: %v", err)
	}
	// new token should validate
	claims, err := svc.ValidateToken(newToken)
	if err != nil {
		t.Fatalf("validate refreshed token: %v", err)
	}
	assert.Equal(t, "user-refresh", claims.UserID)

	// refreshing an invalid token should fail
	_, err = svc.RefreshToken("not-a-token")
	assert.ErrorIs(t, err, utils.ErrInvalidToken)
}
