package http

import (
	"net/http"

	"order-service/domain"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"strconv"
)

type OrderHandler struct {
    usecase *usecase.OrderUsecase
}

func NewOrderHandler(r *gin.Engine, uc *usecase.OrderUsecase) {
    h := &OrderHandler{usecase: uc}
    r.POST("/orders", h.CreateOrder)
	r.GET("/orders/:id", h.GetOrderByID)
	r.PATCH("/orders/:id", h.UpdateOrderStatus)
	r.GET("/orders", h.GetOrdersByUser)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
    var order domain.Order
    if err := c.ShouldBindJSON(&order); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userIDStr, err := c.Cookie("user_id")
    if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
	userID, _ := strconv.Atoi(userIDStr) 
    order.UserID = userID

    if err := h.usecase.CreateOrder(&order); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	order, err := h.usecase.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.UpdateOrderStatus(id, input.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	userStr, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	userID, _ := strconv.Atoi(userStr) 
	orders, err := h.usecase.GetOrdersByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}
