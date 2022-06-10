package handler

import (
	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type User interface {
	CreateUser(user domain.User) (string, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (string, error)
}

type Categories interface {
	Create(list domain.CategoriesList) (string, error)
	GetAll() ([]domain.CategoriesList, error)
	GetById(listId string) (domain.CategoriesList, error)
	Update(itemId string, input domain.UpdateCategoryInput) error
	Delete(itemId string) error
}

type Products interface {
	Create(list domain.ProductsList) (string, error)
	GetAll() ([]domain.ProductsList, error)
	GetById(listId string) (domain.ProductsList, error)
	Update(itemId string, input domain.UpdateProductInput) error
	Delete(itemId string) error
}

type Handler struct {
	userService       User
	categoriesService Categories
	productsService   Products
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		userService:       services.User,
		categoriesService: services.CategoriesList,
		productsService:   services.ProductsList,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.Use(h.loggingMiddleware)

		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api")
	{
		api.Use(h.loggingMiddleware)
		// api.Use(h.userIdentify)

		categories := api.Group("/categories")
		{
			categories.POST("/", h.createCategory)
			categories.GET("/", h.getAllCategories)
			categories.GET("/:id", h.getCategoryById)
			categories.PUT("/:id", h.updateCategory)
			categories.DELETE("/:id", h.deleteCategory)
		}

		products := api.Group("/products")
		{
			products.POST("/", h.createProduct)
			products.GET("/", h.getAllProducts)
			products.GET("/:id", h.getProductById)
			products.PUT("/:id", h.updateProduct)
			products.DELETE("/:id", h.deleteProduct)
		}
	}

	return router
}
