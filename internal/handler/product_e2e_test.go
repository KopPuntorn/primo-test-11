package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"primo-11/internal/database"
	"primo-11/internal/domain"
	"primo-11/internal/handler"
	"primo-11/internal/repository"
	"primo-11/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func setupE2E(t *testing.T) (*echo.Echo, *gorm.DB) {

	host := os.Getenv("DB_HOST")
	if host == "" {
		t.Skip("Skipping E2E test: DB_HOST not set")
	}

	db, err := database.NewPostgresDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	db.Exec("TRUNCATE TABLE products RESTART IDENTITY")

	repo := repository.NewProductRepository(db)
	uc := usecase.NewProductUseCase(repo)
	h := handler.NewProductHandler(uc)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	h.RegisterRoutes(e)

	return e, db
}

func TestProduct_E2E_Flow(t *testing.T) {
	e, db := setupE2E(t)

	createReq := map[string]interface{}{
		"name":        "E2E Product",
		"description": "E2E Description",
		"price":       100.0,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var createResp domain.CreateProductResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createResp)
	assert.NoError(t, err)
	assert.True(t, createResp.Successful)

	var product domain.Product
	db.First(&product)
	assert.Equal(t, "E2E Product", product.Name)
	assert.Equal(t, "E2E Description", *product.Description)
	assert.Equal(t, 100.0, product.Price)

	patchReq := `{"description": null, "price": 150.0}`
	req = httptest.NewRequest(http.MethodPatch, "/product/1", bytes.NewBufferString(patchReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var patchedProduct domain.Product
	db.First(&patchedProduct, 1)
	assert.Equal(t, 150.0, patchedProduct.Price)
	assert.Nil(t, patchedProduct.Description)
	assert.Equal(t, "E2E Product", patchedProduct.Name) 

	
	patchReq2 := `{"name": "Updated Name"}`
	req = httptest.NewRequest(http.MethodPatch, "/product/1", bytes.NewBufferString(patchReq2))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var patchedProduct2 domain.Product
	db.First(&patchedProduct2, 1)
	assert.Equal(t, "Updated Name", patchedProduct2.Name)
	assert.Nil(t, patchedProduct2.Description)   
	assert.Equal(t, 150.0, patchedProduct2.Price) 

	failPatchReq := `{"sale_price": 200.0}`
	req = httptest.NewRequest(http.MethodPatch, "/product/1", bytes.NewBufferString(failPatchReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
