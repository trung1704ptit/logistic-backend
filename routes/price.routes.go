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
	router := rg.Group("contractors/:contractorId/pricing")

	// Route to create new pricing
	router.POST("", rc.pricingController.CreatePricing)

	// Route to get all pricing entries for a contractor
	router.GET("", rc.pricingController.FindAllPricing)

	// Route to get the current pricing for a contractor
	router.GET("/current", rc.pricingController.GetCurrentPricing)

	// Route to update a pricing entry by its ID
	router.PUT("/:pricingId", rc.pricingController.UpdatePricing)

	// Route to delete all pricing entries for a contractor
	router.DELETE("", rc.pricingController.DeletePricing)
}
