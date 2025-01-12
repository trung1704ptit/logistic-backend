package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
	"github.com/wpcodevo/golang-gorm-postgres/middleware"
)

type OrderRouteController struct {
	orderController controllers.OrderController
}

func NewOrderRouteController(orderController controllers.OrderController) OrderRouteController {
	return OrderRouteController{orderController}
}

func (rc *OrderRouteController) OrderRoute(rg *gin.RouterGroup) {
	router := rg.Group("orders")
	router.Use(middleware.DeserializeUser())

	router.POST("", rc.orderController.CreateOrder)            // Create a new order
	router.GET("", rc.orderController.GetOrders)               // Get all orders
	router.GET("/:orderId", rc.orderController.GetOrderByID)   // Get a specific order by ID
	router.PUT("/:orderId", rc.orderController.UpdateOrder)    // Update an order by ID
	router.DELETE("/:orderId", rc.orderController.DeleteOrder) // Delete an order by ID
}
