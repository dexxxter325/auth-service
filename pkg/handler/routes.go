package handler

import (
	_ "CRUD_API/docs"
	"CRUD_API/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine { //инициализация маршрутов
	r := gin.New()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiGroup := r.Group("/api", h.userIdentity)
	apiGroup.POST("/product", h.CreateProduct)
	apiGroup.GET("/product", h.ReadAllProducts)
	apiGroup.GET("/product/:id", h.ReadProductById)
	apiGroup.PUT("/product/:id", h.UpdateProduct)
	apiGroup.DELETE("/product/:id", h.DeleteProduct)

	authGroup := r.Group("/auth")
	authGroup.POST("/sign-up", h.signUp)
	authGroup.POST("/sign-in", h.signIn)
	authGroup.POST("/refresh-tokens", h.Refresh)
	return r
}
