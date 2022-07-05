package handler

import (
	_ "github.com/AndrewMislyuk/go-shop-backend/docs"
	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type User interface {
	CreateUser(user domain.UserSignUp) (string, error)
	GenerateToken(email, password string) (string, error)
	GetMe(token string) (domain.User, error)
}

type Products interface {
	Create(list domain.CreateProductInput) (string, error)
	GetAll() ([]domain.ProductsList, error)
	GetById(listId string) (domain.ProductsList, error)
	Update(itemId string, input domain.UpdateProductInput) error
	Delete(itemId string) error
}

type Files interface {
	Upload(file domain.File) (string, error)
}

type Handler struct {
	userService     User
	productsService Products
	fileService     Files
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		userService:     services.User,
		productsService: services.ProductsList,
		fileService:     services.Files,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(h.CORSMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	auth := router.Group("/auth")
	{
		auth.Use(h.loggingMiddleware)

		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.GET("/get-me", h.userIdentify, h.getMe)
	}

	api := router.Group("/api")
	{
		api.Use(h.loggingMiddleware)

		products := api.Group("/products")
		{
			products.POST("/", h.userIdentify, h.userIsAdmin, h.createProduct)
			products.GET("/", h.getAllProducts)
			products.GET("/:id", h.getProductById)
			products.PUT("/:id", h.userIdentify, h.userIsAdmin, h.updateProduct)
			products.DELETE("/:id", h.userIdentify, h.userIsAdmin, h.deleteProduct)
		}

		files := api.Group("/file")
		{
			files.POST("/upload", h.userIdentify, h.userIsAdmin, h.uploadImage)
		}
	}

	return router
}
