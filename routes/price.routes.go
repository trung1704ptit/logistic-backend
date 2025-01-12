package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
	"github.com/wpcodevo/golang-gorm-postgres/middleware"
)

type PricingRouteController struct {
	pricingController controllers.PricingController
}

func NewPricingRouteController(pricingController controllers.PricingController) PricingRouteController {
	return PricingRouteController{pricingController}
}

func (rc *PricingRouteController) PricingRoute(rg *gin.RouterGroup) {
	router := rg.Group("prices")
	router.Use(middleware.DeserializeUser())

	// Route to create new pricing
	router.POST("/:contractorId", rc.pricingController.CreatePricing)

	router.GET("/:contractorId", rc.pricingController.FindPricingListByContractorID)

	router.GET("/:contractorId/latest", rc.pricingController.FindLatestPricingByContractorID)

	router.DELETE("/:contractorId/:priceId", rc.pricingController.DeletePricingWithDetails)

	router.DELETE("/:contractorId", rc.pricingController.DeleteAllPricingByContractorID)
}
