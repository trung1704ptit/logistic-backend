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
func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.ID = uuid.New() // Generate a new UUID for the order
	if err := ctrl.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrders retrieves all orders
func (ctrl *OrderController) GetOrders(c *gin.Context) {
	var orders []models.Order
	if err := ctrl.DB.Preload("Contractor").
		Preload("Driver").
		Preload("Truck").
		Preload("OtherFees").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder retrieves a specific order by ID
func (ctrl *OrderController) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := ctrl.DB.Preload("Contractor").
		Preload("Driver").
		Preload("Truck").
		Preload("OtherFees").
		First(&order, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order"})
		}
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrder updates an existing order
func (ctrl *OrderController) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := ctrl.DB.First(&order, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order"})
		}
		return
	}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// DeleteOrder deletes an order by ID
func (ctrl *OrderController) DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.DB.Delete(&models.Order{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
