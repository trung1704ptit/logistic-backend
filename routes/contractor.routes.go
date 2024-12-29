package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
)

type ContractorRouteController struct {
	contractorController controllers.ContractorController
}

func NewContractorRouteController(contractorController controllers.ContractorController) ContractorRouteController {
	return ContractorRouteController{contractorController}
}

func (rc *ContractorRouteController) ContractorRoute(rg *gin.RouterGroup) {
	router := rg.Group("contractor")

	router.POST("/", rc.contractorController.CreateContractor)
	router.GET("/", rc.contractorController.FindContractors)
	router.POST("/:contractorId", rc.contractorController.UpdateContractor)
	router.GET("/:contractorId", rc.contractorController.DeleteContractor)
}
