package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-gorm-postgres/models"
	"gorm.io/gorm"
)

type OrderController struct {
	DB *gorm.DB
}

func NewOrderController(DB *gorm.DB) OrderController {
	return OrderController{DB}
}

// CreateOrder creates a new order
func (ctrl *OrderController) CreateOrder(ctx *gin.Context) {
	var newOrder models.Order
	if err := ctx.ShouldBindJSON(&newOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newOrder.ID = uuid.New() // Generate a new UUID for the order
	if err := ctrl.DB.Create(&newOrder).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create order"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newOrder})
}

// GetOrders retrieves all orders
func (ctrl *OrderController) GetOrders(ctx *gin.Context) {
	var orders []models.Order
	if err := ctrl.DB.Preload("Contractor").
		Preload("Driver").
		Preload("Truck").
		Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve orders"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": orders})
}

// GetOrder retrieves a specific order by ID
func (ctrl *OrderController) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := ctrl.DB.Preload("Contractor").
		Preload("Driver").
		Preload("Truck").
		First(&order, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve order"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": order})
}

// UpdateOrder updates an existing order
func (ctrl *OrderController) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := ctrl.DB.First(&order, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve order"})
		}
		return
	}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := ctrl.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": order})
}

// DeleteOrder deletes an order by ID
func (ctrl *OrderController) DeleteOrder(ctx *gin.Context) {
	id := ctx.Param("orderId")

	if err := ctrl.DB.Delete(&models.Order{}, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete order"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
