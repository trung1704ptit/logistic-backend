package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
	"github.com/wpcodevo/golang-gorm-postgres/middleware"
)

type ContractorRouteController struct {
	contractorController controllers.ContractorController
}

func NewContractorRouteController(contractorController controllers.ContractorController) ContractorRouteController {
	return ContractorRouteController{contractorController}
}

func (rc *ContractorRouteController) ContractorRoute(rg *gin.RouterGroup) {
	router := rg.Group("contractors")
	router.Use(middleware.DeserializeUser())

	router.POST("", rc.contractorController.CreateContractor)
	router.GET("", rc.contractorController.FindContractors)
	router.PUT("/:contractorId", rc.contractorController.UpdateContractor)
	router.DELETE("/:contractorId", rc.contractorController.DeleteContractor)
}
