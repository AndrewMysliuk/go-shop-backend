package handler

import (
	"net/http"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type getAllCategoriesListsResponse struct {
	Data []domain.CategoriesList `json:"data"`
}

type getCategoryResponse struct {
	Data domain.CategoriesList `json:"data"`
}

// @Summary Create Category
// @Security ApiKeyAuth
// @Tags Category
// @Description create category
// @ID create-category
// @Accept  json
// @Produce  json
// @Param input body domain.CreateCategoryInput true "Category info"
// @Success 200 {object} getCreationId
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/categories/ [post]
func (h *Handler) createCategory(c *gin.Context) {
	var input domain.CreateCategoryInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.categoriesService.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getCreationId{
		Id: id,
	})
}

// @Summary Get Categories List
// @Tags Category
// @Description get categories list
// @ID get-categories
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllCategoriesListsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/categories/ [get]
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

// @Summary Get Category By ID
// @Tags Category
// @Description get category by id
// @ID get-category-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} getCategoryResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/categories/{id} [get]
func (h *Handler) getCategoryById(c *gin.Context) {
	category_id := c.Param("id")

	category, err := h.categoriesService.GetById(category_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getCategoryResponse{
		Data: category,
	})
}

// @Summary Update Category
// @Security ApiKeyAuth
// @Tags Category
// @Description update category by id
// @ID update-category-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Param input body domain.UpdateCategoryInput true "Category info"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/categories/{id} [put]
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

// @Summary Delete Category
// @Security ApiKeyAuth
// @Tags Category
// @Description delete category by id
// @ID delete-category-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/categories/{id} [delete]
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
