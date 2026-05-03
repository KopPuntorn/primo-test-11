package handler

import (
	"errors"
	"net/http"
	"strconv"

	"primo-11/internal/domain"
	"primo-11/internal/usecase"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	uc usecase.ProductUseCase
}

func NewProductHandler(uc usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

func (h *ProductHandler) RegisterRoutes(e *echo.Echo) {
	product := e.Group("/product")
	product.POST("", h.CreateProduct)
	product.PATCH("/:id", h.PatchProduct)
}

// @Summary      Create Product
// @Description  Create a new product.
// @Accept       json
// @Produce      json
// @Param        product  body      domain.CreateProductRequest  true  "Product details"
// @Success      200      {object}  domain.CreateProductResponse{data=domain.CreateProductData}
// @Router       /product [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req domain.CreateProductRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.CreateProductResponse{
			Successful: false,
			ErrorCode:  "invalid input",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.CreateProductResponse{
			Successful: false,
			ErrorCode:  "invalid input",
		})
	}

	product, err := h.uc.CreateProduct(&req)
	if err != nil {
		status, code := mapError(err)
		return c.JSON(status, domain.CreateProductResponse{
			Successful: false,
			ErrorCode:  code,
		})
	}

	var data2 string
	if product.Description != nil {
		data2 = *product.Description
	}

	return c.JSON(http.StatusOK, domain.CreateProductResponse{
		Successful: true,
		ErrorCode:  "",
		Data: domain.CreateProductData{
			Data1: product.Name,
			Data2: data2,
		},
	})
}

// @Summary      Update Product
// @Description  Update a product's details by ID.
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "Product ID"
// @Param        product  body      domain.PatchProductRequest   true  "Product fields to update"
// @Success      200      {object}  domain.PatchProductResponse
// @Router       /product/{id} [patch]
func (h *ProductHandler) PatchProduct(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.PatchProductResponse{
			Successful: false,
			ErrorCode:  "invalid id",
		})
	}

	var raw map[string]interface{}
	if err := c.Bind(&raw); err != nil {
		return c.JSON(http.StatusBadRequest, domain.PatchProductResponse{
			Successful: false,
			ErrorCode:  "invalid input",
		})
	}

	var req domain.PatchProductRequest
	if err := h.uc.PatchProduct(id, &req, raw); err != nil {
		status, code := mapError(err)
		return c.JSON(status, domain.PatchProductResponse{
			Successful: false,
			ErrorCode:  code,
		})
	}

	return c.JSON(http.StatusOK, domain.PatchProductResponse{
		Successful: true,
		ErrorCode:  "",
	})
}

func parseID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(id), err
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		return http.StatusNotFound, "product not found"
	case errors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
