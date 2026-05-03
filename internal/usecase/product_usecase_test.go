package usecase_test

import (
	"errors"
	"testing"

	"primo-11/internal/domain"
	"primo-11/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) FindByID(id uint) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Update(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func strPtr(s string) *string    { return &s }
func f64Ptr(f float64) *float64  { return &f }
func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUseCase(mockRepo)

	req := &domain.CreateProductRequest{
		Name:  "Test Product",
		Price: 100.0,
	}
	mockRepo.On("Create", mock.AnythingOfType("*domain.Product")).Return(nil)

	product, err := uc.CreateProduct(req)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, 100.0, product.Price)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_WithSalePrice_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUseCase(mockRepo)

	req := &domain.CreateProductRequest{
		Name:      "Sale Item",
		Price:     200.0,
		SalePrice: f64Ptr(150.0),
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Product")).Return(nil)

	product, err := uc.CreateProduct(req)

	assert.NoError(t, err)
	assert.NotNil(t, product.SalePrice)
	assert.Equal(t, 150.0, *product.SalePrice)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_SalePriceHigherThanPrice_ReturnsError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUseCase(mockRepo)

	req := &domain.CreateProductRequest{
		Name:      "Bad Product",
		Price:     100.0,
		SalePrice: f64Ptr(150.0), 
	}
	product, err := uc.CreateProduct(req)
	assert.Nil(t, product)
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateProduct_RepositoryError_ReturnsInternalError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUseCase(mockRepo)

	req := &domain.CreateProductRequest{Name: "Test", Price: 100.0}
	mockRepo.On("Create", mock.AnythingOfType("*domain.Product")).Return(errors.New("db error"))

	product, err := uc.CreateProduct(req)

	assert.Nil(t, product)
	assert.ErrorIs(t, err, domain.ErrInternalServer)
	mockRepo.AssertExpectations(t)
}

func TestPatchProduct_NameOnly_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUseCase(mockRepo)

	existing := &domain.Product{ID: 1, Name: "Old Name", Price: 100.0}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Product")).Return(nil)

	req := &domain.PatchProductRequest{Name: strPtr("Updated")}
	raw := map[string]interface{}{"name": "Updated"}

	err := uc.PatchProduct(1, req, raw)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", existing.Name)
	assert.Equal(t, 100.0, existing.Price)
	mockRepo.AssertExpectations(t)
}

func TestPatchProduct_ProductNotFound_ReturnsError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUseCase(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := uc.PatchProduct(999, &domain.PatchProductRequest{}, nil)

	assert.ErrorIs(t, err, domain.ErrProductNotFound)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestPatchProduct_ClearSalePrice_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := usecase.NewProductUseCase(mockRepo)
	salePrice := 80.0
	existing := &domain.Product{ID: 1, Name: "Item", Price: 100.0, SalePrice: &salePrice}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Product")).Return(nil)

	req := &domain.PatchProductRequest{}
	raw := map[string]interface{}{"sale_price": nil} // Set to null

	err := uc.PatchProduct(1, req, raw)
	assert.NoError(t, err)
	assert.Nil(t, existing.SalePrice)
	mockRepo.AssertExpectations(t)
}
