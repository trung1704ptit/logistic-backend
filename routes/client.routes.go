package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
	"github.com/wpcodevo/golang-gorm-postgres/middleware"
)

type ClientRouteController struct {
	ClientController controllers.ClientController
}

func NewClientRouteController(ClientController controllers.ClientController) ClientRouteController {
	return ClientRouteController{ClientController}
}

func (rc *ClientRouteController) ClientRoute(rg *gin.RouterGroup) {
	router := rg.Group("clients")
	router.Use(middleware.DeserializeUser())

	router.POST("", rc.ClientController.CreateClient)
	router.GET("", rc.ClientController.FindClients)
	router.PUT("/:clientId", rc.ClientController.UpdateClient)
	router.DELETE("/:clientId", rc.ClientController.DeleteClient)
}
