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
	router.POST("/:ownerId", rc.pricingController.CreatePricing)

	router.GET("/:ownerId", rc.pricingController.FindPricingListByOwner)

	router.GET("/:ownerId/latest", rc.pricingController.FindLatestPricingByOwner)

	router.GET("/:ownerId/:priceId", rc.pricingController.FindPricingByOwnerAndPriceID)

	router.DELETE("/:ownerId/:priceId", rc.pricingController.DeletePricingWithDetails)

	router.DELETE("/:ownerId", rc.pricingController.DeleteAllPricingByContractorID)
}
