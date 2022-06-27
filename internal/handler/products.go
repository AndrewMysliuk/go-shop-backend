package handler

import (
	"net/http"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type getAllProductsListsResponse struct {
	Data []domain.ProductsList `json:"data"`
}

type getProductResponse struct {
	Data domain.ProductsList `json:"data"`
}

// @Summary Create Product
// @Security ApiKeyAuth
// @Tags Product
// @Description create product
// @ID create-product
// @Accept  json
// @Produce  json
// @Param input body domain.CreateProductInput true "Product info"
// @Success 200 {object} getCreationId
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/products/ [post]
func (h *Handler) createProduct(c *gin.Context) {
	var input domain.CreateProductInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	id, err := h.productsService.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, getCreationId{
		Id: id,
	})
}

// @Summary Get Products List
// @Tags Product
// @Description get products list
// @ID get-products
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllProductsListsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/products/ [get]
func (h *Handler) getAllProducts(c *gin.Context) {
	products, err := h.productsService.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, getAllProductsListsResponse{
		Data: products,
	})
}

// @Summary Get Product By ID
// @Tags Product
// @Description get product by id
// @ID get-product-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} getProductResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/products/{id} [get]
func (h *Handler) getProductById(c *gin.Context) {
	product_id := c.Param("id")

	product, err := h.productsService.GetById(product_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, getProductResponse{
		Data: product,
	})
}

// @Summary Update Product
// @Security ApiKeyAuth
// @Tags Product
// @Description update product by id
// @ID update-product-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param input body domain.UpdateProductInput true "Product info"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/products/{id} [put]
func (h *Handler) updateProduct(c *gin.Context) {
	product_id := c.Param("id")

	var input domain.UpdateProductInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.productsService.Update(product_id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary Delete Product
// @Security ApiKeyAuth
// @Tags Product
// @Description delete product by id
// @ID delete-product-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/products/{id} [delete]
func (h *Handler) deleteProduct(c *gin.Context) {
	product_id := c.Param("id")

	err := h.productsService.Delete(product_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
