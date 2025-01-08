package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
)

type PricingRouteController struct {
	pricingController controllers.PricingController
}

func NewPricingRouteController(pricingController controllers.PricingController) PricingRouteController {
	return PricingRouteController{pricingController}
}

func (rc *PricingRouteController) PricingRoute(rg *gin.RouterGroup) {
	router := rg.Group("prices")

	// Route to create new pricing
	router.POST("/:contractorId", rc.pricingController.CreatePricing)

	// // Route to get all pricing entries for a contractor
	router.GET("/:contractorId", rc.pricingController.GetPricingByContractorID)

	// // Route to get the current pricing for a contractor
	router.GET("/:contractorId/current", rc.pricingController.GetLatestPricingByContractorID)

	// // Route to update a pricing entry by its ID
	// router.PUT("/:contractorId/:pricingId", rc.pricingController.UpdatePricing)

	// // Route to delete all pricing entries for a contractor
	// router.DELETE("/:contractorId", rc.pricingController.DeletePricing)
}
