package domain

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"not null"                 json:"name"`
	Description *string   `                                json:"description"`
	Price       float64   `gorm:"not null"                 json:"price"`
	SalePrice   *float64  `                                json:"sale_price"`
	CreatedAt   time.Time `                                json:"created_at"`
	UpdatedAt   time.Time `                                json:"updated_at"`
}

type ProductRepository interface {
	Create(product *Product) error
	Update(product *Product) error
	FindByID(id uint) (*Product, error)
}

type CreateProductRequest struct {
	Name        string   `json:"name"        validate:"required"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price"       validate:"required,gt=0"`
}

type PatchProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       *float64 `json:"price"`
}

type CreateProductResponse struct {
	Successful bool        `json:"successful"`
	ErrorCode  string      `json:"error_code"`
	Data       interface{} `json:"data"`
}

// Data สำหรับ CreateProduct ตามโจทย์ (data1, data2)
type CreateProductData struct {
	Data1 string `json:"data1"`
	Data2 string `json:"data2"`
}

type PatchProductResponse struct {
	Successful bool   `json:"successful"`
	ErrorCode  string `json:"error_code"`
}
