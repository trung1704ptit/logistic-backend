package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
)

type TruckRouteController struct {
	truckController controllers.TruckController
}

func NewTruckRouteController(truckController controllers.TruckController) TruckRouteController {
	return TruckRouteController{truckController}
}

func (rc *TruckRouteController) TruckRoute(rg *gin.RouterGroup) {
	router := rg.Group("truck")

	router.POST("/", rc.truckController.CreateTruck)
	router.GET("/", rc.truckController.FindTrucks)
	router.POST("/:truckId", rc.truckController.UpdateTruck)
	router.GET("/:truckId", rc.truckController.FindTruckById)
	router.DELETE("/:truckId", rc.truckController.DeleteTruck)
}
