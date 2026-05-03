package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"primo-11/internal/domain"
	"primo-11/internal/handler"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductUseCase struct {
	mock.Mock
}

func (m *MockProductUseCase) CreateProduct(req *domain.CreateProductRequest) (*domain.Product, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductUseCase) PatchProduct(id uint, req *domain.PatchProductRequest, raw map[string]interface{}) error {
	args := m.Called(id, req, raw)
	return args.Error(0)
}

func (m *MockProductUseCase) GetProduct(id uint) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductUseCase) ListProducts() ([]*domain.Product, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *MockProductUseCase) DeleteProduct(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupEcho(uc *MockProductUseCase) *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	h := handler.NewProductHandler(uc)
	h.RegisterRoutes(e)
	return e
}

func makeRequest(method, path string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return req, rec
}

func TestCreateProductHandler_Success(t *testing.T) {
	mockUC := new(MockProductUseCase)
	e := setupEcho(mockUC)

	desc := "A great product"
	expectedProduct := &domain.Product{
		ID:          1,
		Name:        "Test Product",
		Description: &desc,
		Price:       99.99,
	}

	mockUC.On("CreateProduct", mock.AnythingOfType("*domain.CreateProductRequest")).
		Return(expectedProduct, nil)

	req, rec := makeRequest(http.MethodPost, "/product", map[string]interface{}{
		"name":        "Test Product",
		"description": "A great product",
		"price":       99.99,
	})

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp domain.CreateProductResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.True(t, resp.Successful)
	
	// ตรวจสอบว่า Data มี data1 (โจทย์ใหม่)
	dataMap := resp.Data.(map[string]interface{})
	assert.Equal(t, "Test Product", dataMap["data1"])

	mockUC.AssertExpectations(t)
}

func TestCreateProductHandler_MissingRequiredField(t *testing.T) {
	mockUC := new(MockProductUseCase)
	e := setupEcho(mockUC)

	req, rec := makeRequest(http.MethodPost, "/product", map[string]interface{}{
		"price": 99.99,
	})

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp domain.CreateProductResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.False(t, resp.Successful)
	mockUC.AssertNotCalled(t, "CreateProduct")
}

func TestPatchProductHandler_Success(t *testing.T) {
	mockUC := new(MockProductUseCase)
	e := setupEcho(mockUC)

	mockUC.On("PatchProduct", uint(1), mock.Anything, mock.Anything).Return(nil)

	req, rec := makeRequest(http.MethodPatch, "/product/1", map[string]interface{}{
		"name": "Updated Name",
	})

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp domain.PatchProductResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.True(t, resp.Successful)
	mockUC.AssertExpectations(t)
}

func TestPatchProductHandler_NotFound(t *testing.T) {
	mockUC := new(MockProductUseCase)
	e := setupEcho(mockUC)

	mockUC.On("PatchProduct", uint(999), mock.Anything, mock.Anything).Return(domain.ErrProductNotFound)

	req, rec := makeRequest(http.MethodPatch, "/product/999", map[string]interface{}{
		"name": "X",
	})
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockUC.AssertExpectations(t)
}
