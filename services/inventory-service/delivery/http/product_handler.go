package http

import (
	"net/http"
	"strconv"
    // "context"
	"inventory-service/domain"
	"inventory-service/internal/usecase"

	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson"
)

type ProductHandler struct {
    usecase *usecase.ProductUsecase
}

func NewProductHandler(router *gin.Engine, usecase *usecase.ProductUsecase) {
    handler := &ProductHandler{usecase: usecase}

    router.GET("/products/:id", handler.GetProduct)
	router.POST("/products", handler.CreateProduct)
    // router.PATCH("/products/:id", handler.UpdateProduct)
    router.DELETE("/products/:id", handler.DeleteProduct)
    router.GET("/products", handler.GetAllProducts)

}

func (h *ProductHandler) GetProduct(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    product, err := h.usecase.GetProduct(int32(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }
    c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
    var product domain.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    createdProduct, err := h.usecase.CreateProduct(product)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, createdProduct)

}

// func (h *ProductHandler) UpdateProduct(c *gin.Context) {
// 	id := (c.Param("id"))

// 	var update bson.M
// 	if err := c.ShouldBindJSON(&update); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	err := h.usecase.UpdateProduct(update)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
// }

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.usecase.DeleteProduct(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}


func (h *ProductHandler) GetAllProducts(c *gin.Context) {
    name := c.DefaultQuery("name", "") 
    category, _ := strconv.Atoi(c.DefaultQuery("category", ""))
    limitStr := c.DefaultQuery("page", "10") 
    skipStr := c.DefaultQuery("page_size", "0") 

    limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skip, _ := strconv.ParseInt(skipStr, 10, 64)

    products, err := h.usecase.GetAllProducts(name, category, int(limit), int(skip))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    c.JSON(http.StatusOK, products)
}