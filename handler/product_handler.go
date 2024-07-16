package handler

import (
	"net/http"
	"strconv"

	"product-service/model"
	"product-service/usecase"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	Usecase usecase.IProductUsecase
}

func NewProductHandler(e *echo.Group, u usecase.IProductUsecase) {
	handler := &ProductHandler{
		Usecase: u,
	}

	e.GET("/products", handler.GetProducts)
	e.GET("/products/:id", handler.GetProductByID)
	e.POST("/products", handler.CreateProduct)
	e.PUT("/products/:id", handler.UpdateProduct)
	e.DELETE("/products/:id", handler.DeleteProduct)
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0 // Default offset
	}

	products, err := h.Usecase.GetProducts(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	product, err := h.Usecase.GetProductByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var product model.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.Usecase.CreateProduct(c.Request().Context(), product); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "Product created successfully")
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	var product model.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "Invalid product ID")
	}
	product.ID = id

	if err := h.Usecase.UpdateProduct(c.Request().Context(), product); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Product updated successfully")
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	if err := h.Usecase.DeleteProduct(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Product deleted successfully")
}
