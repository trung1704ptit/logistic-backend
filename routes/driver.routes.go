package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
	"github.com/wpcodevo/golang-gorm-postgres/middleware"
)

type DriverRouteController struct {
	driverController controllers.DriverController
}

func NewDriverRouteController(driverController controllers.DriverController) DriverRouteController {
	return DriverRouteController{driverController}
}

func (rc *DriverRouteController) DriverRoute(rg *gin.RouterGroup) {
	router := rg.Group("drivers")
	router.Use(middleware.DeserializeUser())

	router.POST("", rc.driverController.CreateDriver)
	router.GET("", rc.driverController.FindDrivers)
	router.POST("/:driverId", rc.driverController.UpdateDriver)
	router.GET("/:driverId", rc.driverController.FindDriverById)
	router.PUT("/:driverId", rc.driverController.UpdateDriver)
	router.DELETE("/:driverId", rc.driverController.DeleteDriver)
}
