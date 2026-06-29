//go:build integration
// +build integration

package categories_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nicolas/finanzas/backend/internal/auth"
	"github.com/nicolas/finanzas/backend/internal/categories"
	"github.com/nicolas/finanzas/backend/internal/testhelpers"
)

func setupCategoryRepo(t *testing.T) (categories.CategoryRepository, uuid.UUID, func()) {
	t.Helper()
	db := testhelpers.SetupTestDB(t)
	userRepo := auth.NewUserRepository(db.DB)
	email := uuid.NewString() + "@example.com"
	require.NoError(t, userRepo.Create(&auth.User{Email: email, PasswordHash: "x"}))
	user, err := userRepo.FindByEmail(email)
	require.NoError(t, err)
	return categories.NewCategoryRepository(db.DB), user.ID, db.Cleanup
}

func TestCategoryRepository_CreateAndList(t *testing.T) {
	repo, userID, cleanup := setupCategoryRepo(t)
	defer cleanup()

	c := &categories.Category{
		UserID: userID,
		Name:   "Comida",
		Type:   categories.TypeExpense,
	}
	require.NoError(t, repo.Create(c))

	list, err := repo.ListByUser(userID, "")
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "Comida", list[0].Name)

	// Filter by type
	list, err = repo.ListByUser(userID, "income")
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestCategoryRepository_GetByID_OwnershipEnforced(t *testing.T) {
	repo, userID, cleanup := setupCategoryRepo(t)
	defer cleanup()

	c := &categories.Category{UserID: userID, Name: "X", Type: categories.TypeExpense}
	require.NoError(t, repo.Create(c))

	found, err := repo.GetByID(c.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, "X", found.Name)

	_, err = repo.GetByID(c.ID, uuid.New())
	assert.ErrorIs(t, err, categories.ErrCategoryNotFound)
}

func TestCategoryRepository_SoftDelete(t *testing.T) {
	repo, userID, cleanup := setupCategoryRepo(t)
	defer cleanup()

	c := &categories.Category{UserID: userID, Name: "X", Type: categories.TypeExpense}
	require.NoError(t, repo.Create(c))

	require.NoError(t, repo.Delete(c.ID, userID))

	list, _ := repo.ListByUser(userID, "")
	assert.Empty(t, list)
}

func TestService_SeedDefaults(t *testing.T) {
	db := testhelpers.SetupTestDB(t)
	defer db.Cleanup()
	userRepo := auth.NewUserRepository(db.DB)
	email := uuid.NewString() + "@example.com"
	require.NoError(t, userRepo.Create(&auth.User{Email: email, PasswordHash: "x"}))
	user, _ := userRepo.FindByEmail(email)

	repo := categories.NewCategoryRepository(db.DB)
	svc := categories.NewService(repo)

	created, err := svc.SeedDefaults(user.ID)
	require.NoError(t, err)
	assert.Equal(t, len(categories.DefaultCategoriesES), created)

	// Idempotent
	created2, err := svc.SeedDefaults(user.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, created2)

	list, _ := repo.ListByUser(user.ID, "")
	assert.Len(t, list, len(categories.DefaultCategoriesES))
}