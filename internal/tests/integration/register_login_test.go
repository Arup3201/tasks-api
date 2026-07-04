package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Arup3201/gotask/internal/config"
	"github.com/Arup3201/gotask/internal/models"
	"github.com/Arup3201/gotask/internal/testutils"
	"github.com/Arup3201/gotask/internal/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRegistrationAndLoginWorkflow(t *testing.T) {
	ctx := context.Background()

	pg, err := testutils.CreatePostgresContainer(ctx)
	if err != nil {
		t.Fatalf("could not start postgres container: %v", err)
	}
	defer func() {
		if err := pg.Terminate(ctx); err != nil {
			t.Fatalf("could not terminate postgres container: %v", err)
		}
	}()

	db, err := gorm.Open(postgres.Open(pg.ConnectionString), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate user schema: %v", err)
	}

	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_ISSUER", "test-issuer")
	config := config.Load()

	store := models.NewUserStore(db)
	userService := models.NewUserService(store)
	jwtService := utils.NewJWTService(config)

	createdUser, err := userService.CreateUser(ctx, "alice@example.com", "Alice Example", "secret123")
	if err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	assert.Equal(t, "alice@example.com", createdUser.Email)
	assert.Equal(t, "Alice Example", createdUser.Name)
	assert.NotEmpty(t, createdUser.ID)
	assert.False(t, createdUser.CreatedAt.IsZero())
	assert.False(t, createdUser.UpdatedAt.IsZero())

	loggedInUser, err := userService.ExchangeUserWithCredentials(ctx, "alice@example.com", "secret123")
	if err != nil {
		t.Fatalf("exchange credentials failed: %v", err)
	}

	assert.Equal(t, createdUser.ID, loggedInUser.ID)
	assert.Equal(t, createdUser.Email, loggedInUser.Email)

	token, err := jwtService.GenerateToken(loggedInUser.ID, loggedInUser.Email)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}
	assert.NotEmpty(t, token)

	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		t.Fatalf("validate token failed: %v", err)
	}

	assert.Equal(t, loggedInUser.ID, claims.UserID)
	assert.Equal(t, loggedInUser.Email, claims.Email)
	assert.Equal(t, "test-issuer", claims.Issuer)
	assert.NotNil(t, claims.ExpiresAt)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), claims.ExpiresAt.Time, 10*time.Second)
}
