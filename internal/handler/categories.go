package handler

import (
	"net/http"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type getAllCategoriesListsResponse struct {
	Data []domain.CategoriesList `json:"data"`
}

func (h *Handler) createCategory(c *gin.Context) {
	var input domain.CategoriesList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.categoriesService.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllCategories(c *gin.Context) {
	lists, err := h.categoriesService.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllCategoriesListsResponse{
		Data: lists,
	})
}

func (h *Handler) getCategoryById(c *gin.Context) {
	category_id := c.Param("id")

	category, err := h.categoriesService.GetById(category_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *Handler) updateCategory(c *gin.Context) {
	category_id := c.Param("id")

	var input domain.UpdateCategoryInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.categoriesService.Update(category_id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteCategory(c *gin.Context) {
	category_id := c.Param("id")

	err := h.categoriesService.Delete(category_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
