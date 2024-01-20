package handler

import (
	"CRUD_API"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// @Summary CreateProduct
// @Tags product
// @Security ApiKeyAuth
// @Description create product
// @Accept  json
// @Produce  json
// @Param input body CreateProductRequest true "product info"
// @Router /api/product [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	var product CRUD_API.Products
	if err := c.BindJSON(&product); err != nil { //из json запроса данные присваиваем в нашу структуру
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON in CreateProduct"})
		return
	}
	createdProduct, err := h.services.Create(product.Name, product.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "err in service"})
		return
	}
	cacheKey := fmt.Sprintf("product:%d", createdProduct.ID) //ключ нашего кеша и значение на основе созданного product.(в общем-ключ для хэширования)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error marshaling product to JSON": err.Error()})
		return
	}
	// Шаг 3: Сохранение данных в кэш
	//создаем jitter(если много кэшей закончатся одновременно-будет сбой,мы делаем рандомный пинг к завершению кэшей)
	expiration := cacheTTl + time.Duration(rand.Intn(21)-10)*time.Second
	/*рандомное колво сек от 0 до20,затем - 10<>создаем интервал возможных знач от -10с до 10с. Это вычитается или прибавляется к cacheTTl.*/
	cache := GetSingleton()
	err = cache.Set(cacheKey, createdProduct, expiration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in Set data to redis(create)": err.Error()})
		return
	}

	c.JSON(200, gin.H{"product(created successfully)": createdProduct})
}

type CreateProductRequest struct { //for swagger
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// @Summary ReadAllProducts
// @Tags product
// @Security ApiKeyAuth
// @Description readallproducts!
// @Accept json
// @Produce json
// @Router /api/product [get]
func (h *Handler) ReadAllProducts(c *gin.Context) {
	AllProducts, err := h.services.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "err in service"})
		return
	}
	c.JSON(200, gin.H{"product": AllProducts})
}

// @Summary ReadProductById
// @Tags product
// @Security ApiKeyAuth
// @Description readProductById!!
// @Accept json
// @Produce json
// @Param id path int true "write id"
// @Router /api/product/{id} [get]
func (h *Handler) ReadProductById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	cache := GetSingleton()
	// Шаг 1: Попытка получения данных из кэша
	cacheKey := fmt.Sprintf("product:%d", id)
	cachedProduct := cache.Get(cacheKey)
	/*cachedProduct, err := RedisClient.Get(c, cacheKey).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err in get data from redis": err.Error()})
	}*/

	// Если данные найдены в кэше, возвращаем их
	if cachedProduct != nil {
		c.JSON(http.StatusOK, gin.H{"product": cachedProduct, "from_cache": true})
		return
	}

	// Если данных нет в кэше, получаем их из бд
	product, err := h.services.ReadById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "err in service"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error marshaling product to JSON in readbyid": err.Error()})
		return
	}
	expiration := cacheTTl + time.Duration(rand.Intn(21)-10)*time.Second
	// Сохранение данных в кэш
	err = cache.Set(cacheKey, product, expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err in set data to redis": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"product": product, "from_cache": false})
}

// @Summary UpdateProduct
// @Tags product
// @Security ApiKeyAuth
// @Description Update product by ID
// @Accept json
// @Produce json
// @Param name,description body UpdateProductRequest true "product info"
// @Param id path int true "product info"
// @Router /api/product/{id} [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	var product CRUD_API.Products
	if err := c.BindJSON(&product); err != nil { //новые данные из json запроса привязываем к нашей структуре(обновляем)вместо старых
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json in UpdateProduct"})
		return
	}
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	UpdatedProduct, err := h.services.Update(product.Name, product.Description, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error marshaling product to JSON in update": "err in service"})
		return
	}
	cacheKey := fmt.Sprintf("product:%d", id)
	jsonProduct, err := json.Marshal(UpdatedProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//expiration := cacheTTl + time.Duration(rand.Intn(21)-10)*time.Second
	err = RedisClient.Set(c, cacheKey, jsonProduct, cacheTTl).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in set data to redis(update)": err.Error()})
		return
	}

	c.JSON(200, gin.H{"product(updated successfully)": UpdatedProduct})
}

type UpdateProductRequest struct { //for swagger
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// @Summary DeleteProduct
// @Tags product
// @Security ApiKeyAuth
// @Description Delete product by ID
// @Accept json
// @Produce json
// @Param id path int true "fill ID!"
// @Router /api/product/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	err = h.services.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "err in service"})
		return
	}

	cacheKey := fmt.Sprintf("product:%d", id)
	err = RedisClient.Del(c, cacheKey).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err in delete data from redis": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "deleted successfully")
}
