package identity_test

import (
	"context"
	"testing"
	"time"

	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/charlesvdv/cirrus/backend/database/sqlite"
	"github.com/charlesvdv/cirrus/backend/identity"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repository := sqlite.IdentityRepository{}
	db := sqlite.NewTestDatabase()
	defer db.Close()

	userService := identity.NewUserService(db, repository)

	t.Run("CreateUser", func(t *testing.T) {
		integrationTestCreateUser(t, userService)
	})
}

func integrationTestCreateUser(t *testing.T, service identity.UserService) {
	var user cirrus.User
	err := service.CreateUser(context.Background(), &user)
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)
	require.WithinDuration(t, time.Now(), user.CreatedAt, 5*time.Second)
}
