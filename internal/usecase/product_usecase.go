package usecase

import "primo-11/internal/domain"

type ProductUseCase interface {
	CreateProduct(req *domain.CreateProductRequest) (*domain.Product, error)
	PatchProduct(id uint, req *domain.PatchProductRequest, raw map[string]interface{}) error
}

type productUseCase struct {
	repo domain.ProductRepository
}

func NewProductUseCase(repo domain.ProductRepository) ProductUseCase {
	return &productUseCase{repo: repo}
}

func (uc *productUseCase) CreateProduct(req *domain.CreateProductRequest) (*domain.Product, error) {

	if req.SalePrice != nil && *req.SalePrice >= req.Price {
		return nil, domain.ErrInvalidInput
	}
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SalePrice:   req.SalePrice,
	}

	if err := uc.repo.Create(product); err != nil {
		return nil, domain.ErrInternalServer
	}

	return product, nil
}

func (uc *productUseCase) PatchProduct(id uint, req *domain.PatchProductRequest, raw map[string]interface{}) error {

	product, err := uc.repo.FindByID(id)
	if err != nil {
		return domain.ErrProductNotFound
	}

	if val, ok := raw["name"]; ok && val != nil {
		if s, ok := val.(string); ok {
			product.Name = s
		}
	}

	if val, ok := raw["description"]; ok {
		if val == nil {
			product.Description = nil 
		} else if s, ok := val.(string); ok {
			product.Description = &s
		}
	}

	if val, ok := raw["price"]; ok && val != nil {
		if f, ok := val.(float64); ok {
			product.Price = f
		}
	}

	if val, ok := raw["sale_price"]; ok {
		if val == nil {
			product.SalePrice = nil 
		} else if f, ok := val.(float64); ok {
			product.SalePrice = &f
		}
	}

	if product.SalePrice != nil && *product.SalePrice >= product.Price {
		return domain.ErrInvalidInput
	}

	return uc.repo.Update(product)
}
