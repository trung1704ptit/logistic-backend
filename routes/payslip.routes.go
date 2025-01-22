package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
	"github.com/wpcodevo/golang-gorm-postgres/middleware"
)

type PayslipRouteController struct {
	payslipController controllers.PayslipController
}

func NewPayslipRouteController(payslipController controllers.PayslipController) PayslipRouteController {
	return PayslipRouteController{payslipController}
}

func (rc *PayslipRouteController) PayslipRoute(rg *gin.RouterGroup) {
	router := rg.Group("payslips")
	router.Use(middleware.DeserializeUser())

	router.POST("", rc.payslipController.CreatePayslip)              // Create a new payslip
	router.GET("", rc.payslipController.GetPayslips)                 // Get all payslips
	router.GET("/:payslipId", rc.payslipController.GetPayslipByID)   // Get a specific payslip by ID
	router.PUT("/:payslipId", rc.payslipController.UpdatePayslip)    // Update a payslip by ID
	router.DELETE("/:payslipId", rc.payslipController.DeletePayslip) // Delete a payslip by ID
}
