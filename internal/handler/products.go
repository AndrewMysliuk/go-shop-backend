package handler

import (
	"net/http"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type getAllProductsListsResponse struct {
	Data []domain.ProductsList `json:"data"`
}

func (h *Handler) createProduct(c *gin.Context) {
	var input domain.ProductsList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.productsService.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

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

func (h *Handler) getProductById(c *gin.Context) {
	product_id := c.Param("id")

	product, err := h.productsService.GetById(product_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, product)
}

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
