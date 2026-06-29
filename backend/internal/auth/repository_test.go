//go:build integration
// +build integration

package auth_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nicolas/finanzas/backend/internal/auth"
	"github.com/nicolas/finanzas/backend/internal/testhelpers"
)

func newEmail(t *testing.T) string {
	t.Helper()
	return "user-" + uuid.NewString() + "@example.com"
}

func TestCreateUser_Success(t *testing.T) {
	db := testhelpers.SetupTestDB(t)
	defer db.Cleanup()

	repo := auth.NewUserRepository(db.DB)
	email := newEmail(t)
	user := &auth.User{
		Email:        email,
		PasswordHash: "hashed_password",
	}

	err := repo.Create(user)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.Equal(t, email, user.Email)
}

func TestCreateUser_DuplicateEmail_ReturnsErrUserAlreadyExists(t *testing.T) {
	db := testhelpers.SetupTestDB(t)
	defer db.Cleanup()

	repo := auth.NewUserRepository(db.DB)
	email := newEmail(t)

	first := &auth.User{Email: email, PasswordHash: "hash1"}
	err := repo.Create(first)
	require.NoError(t, err)

	second := &auth.User{Email: email, PasswordHash: "hash2"}
	err = repo.Create(second)
	assert.ErrorIs(t, err, auth.ErrUserAlreadyExists)
}

func TestFindByEmail_Success(t *testing.T) {
	db := testhelpers.SetupTestDB(t)
	defer db.Cleanup()

	repo := auth.NewUserRepository(db.DB)
	email := newEmail(t)
	created := &auth.User{Email: email, PasswordHash: "hash"}
	require.NoError(t, repo.Create(created))

	found, err := repo.FindByEmail(email)
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, email, found.Email)
}

func TestFindByEmail_NotFound(t *testing.T) {
	db := testhelpers.SetupTestDB(t)
	defer db.Cleanup()

	repo := auth.NewUserRepository(db.DB)

	_, err := repo.FindByEmail("nonexistent@example.com")
	assert.ErrorIs(t, err, auth.ErrUserNotFound)
}

func TestFindByID_Success(t *testing.T) {
	db := testhelpers.SetupTestDB(t)
	defer db.Cleanup()

	repo := auth.NewUserRepository(db.DB)
	email := newEmail(t)
	created := &auth.User{Email: email, PasswordHash: "hash"}
	require.NoError(t, repo.Create(created))

	found, err := repo.FindByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, email, found.Email)
}

func TestFindByID_NotFound(t *testing.T) {
	db := testhelpers.SetupTestDB(t)
	defer db.Cleanup()

	repo := auth.NewUserRepository(db.DB)

	_, err := repo.FindByID(uuid.New())
	assert.ErrorIs(t, err, auth.ErrUserNotFound)
}