package controllers

import (
	"net/http"
	"strconv"

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
// func (ctrl *OrderController) GetOrders(ctx *gin.Context) {
// 	var orders []models.Order
// 	if err := ctrl.DB.Preload("Contractor").
// 		Preload("Driver").
// 		Preload("Truck").
// 		Order("updated_at DESC").
// 		Find(&orders).Error; err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve orders"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": orders})
// }

func (ctrl *OrderController) GetOrders(ctx *gin.Context) {
	month := ctx.Query("month")
	year := ctx.Query("year")
	driverId := ctx.Query("driver_id")
	contractorId := ctx.Query("contractor_id")

	// validate the parameters
	if month == "" || year == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Month and Year are required"})
		return
	}

	// convert month and year to integer

	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Invalid month"})
		return
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil || yearInt < 2018 || yearInt > 2100 {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Invalid year"})
		return
	}

	// query database with specific month and year
	var orders []models.Order
	query := ctrl.DB.Preload("Contractor").
		Preload("Driver").
		Preload("Truck").
		Where("EXTRACT(MONTH from updated_at) = ? AND EXTRACT(YEAR FROM updated_at) = ?", monthInt, yearInt)

	// add driver id condition if provided
	if driverId != "all" && driverId != "" {
		query.Where("driver_id = ?", driverId)
	}

	// add contractor id condition if provided
	if contractorId != "all" && contractorId != "" {
		query.Where("contractor_id = ?", contractorId)
	}

	if err := query.Order("updated_at DESC").Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve orders"})
		return
	}
	// Return the orders as a JSON response
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": orders})
}

// GetOrder retrieves a specific order by ID
func (ctrl *OrderController) GetOrderByID(c *gin.Context) {
	id := c.Param("orderId")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid order ID format"})
		return
	}

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
	id := c.Param("orderId")
	var order models.Order

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid order ID format"})
		return
	}

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
