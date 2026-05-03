package repository_test

import (
	"os"
	"testing"

	"primo-11/internal/domain"
	"primo-11/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func setupTestDB(t *testing.T) *gorm.DB {
	host := os.Getenv("DB_HOST")
	if host == "" {
		t.Skip("DB_HOST not set — skipping integration test")
	}

	dsn := "host=" + host +
		" user=" + getEnv("DB_USER", "postgres") +
		" password=" + getEnv("DB_PASSWORD", "postgres") +
		" dbname=" + getEnv("DB_NAME", "product_test") +
		" port=" + getEnv("DB_PORT", "5432") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	db.Migrator().DropTable(&domain.Product{})
	db.AutoMigrate(&domain.Product{})

	return db
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func TestProductRepository_Create_Integration(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProductRepository(db)

	desc := "Integration test product"
	product := &domain.Product{
		Name:        "Integration Test",
		Description: &desc,
		Price:       199.99,
	}

	err := repo.Create(product)

	assert.NoError(t, err)
	assert.NotZero(t, product.ID) 
}

func TestProductRepository_FindByID_Integration(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProductRepository(db)

	
	product := &domain.Product{Name: "Find Me", Price: 50.0}
	require.NoError(t, repo.Create(product))

	found, err := repo.FindByID(product.ID)

	assert.NoError(t, err)
	assert.Equal(t, "Find Me", found.Name)
}

func TestProductRepository_Update_Integration(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProductRepository(db)

	product := &domain.Product{Name: "Before", Price: 100.0}
	require.NoError(t, repo.Create(product))

	product.Name = "After"
	err := repo.Update(product)

	assert.NoError(t, err)

	updated, _ := repo.FindByID(product.ID)
	assert.Equal(t, "After", updated.Name)
}

