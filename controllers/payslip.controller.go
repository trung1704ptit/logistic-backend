package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-gorm-postgres/models"
	"gorm.io/gorm"
)

type PayslipController struct {
	DB *gorm.DB
}

func NewPayslipController(DB *gorm.DB) PayslipController {
	return PayslipController{DB}
}

// CreatePayslip creates a new payslip
func (ctrl *PayslipController) CreatePayslip(ctx *gin.Context) {
	var newPayslip models.Payslip
	if err := ctx.ShouldBindJSON(&newPayslip); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPayslip.ID = uuid.New() // Generate a new UUID for the payslip
	if err := ctrl.DB.Create(&newPayslip).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newPayslip})
}

// GetPayslips retrieves all payslips based on filters
func (ctrl *PayslipController) GetPayslips(ctx *gin.Context) {
	month := ctx.Query("month")
	year := ctx.Query("year")
	driverId := ctx.Query("driver_id")
	contractorId := ctx.Query("contractor_id")

	// Validate the parameters
	if month == "" || year == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Month and Year are required"})
		return
	}

	// Convert month and year to integer
	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid month"})
		return
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil || yearInt < 2018 || yearInt > 2100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid year"})
		return
	}

	// Query database with specific month and year
	var payslips []models.Payslip
	query := ctrl.DB.Preload("Contractor").
		Preload("Driver").
		Where("month = ? AND year = ?", monthInt, yearInt)

	// Add driver ID condition if provided
	if driverId != "all" && driverId != "" {
		query = query.Where("driver_id = ?", driverId)
	}

	// Add contractor ID condition if provided
	if contractorId != "all" && contractorId != "" {
		query = query.Where("contractor_id = ?", contractorId)
	}

	if err := query.Order("updated_at DESC").Find(&payslips).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve payslips"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": payslips})
}

// GetPayslipByID retrieves a specific payslip by ID
func (ctrl *PayslipController) GetPayslipByID(ctx *gin.Context) {
	id := ctx.Param("payslipId")
	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid payslip ID format"})
		return
	}

	var payslip models.Payslip
	if err := ctrl.DB.Preload("Contractor").
		Preload("Driver").
		First(&payslip, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Payslip not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve payslip"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": payslip})
}

// UpdatePayslip updates an existing payslip
func (ctrl *PayslipController) UpdatePayslip(ctx *gin.Context) {
	id := ctx.Param("payslipId")
	var payslip models.Payslip

	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid payslip ID format"})
		return
	}

	if err := ctrl.DB.First(&payslip, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Payslip not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve payslip"})
		}
		return
	}

	if err := ctx.ShouldBindJSON(&payslip); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := ctrl.DB.Save(&payslip).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update payslip"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": payslip})
}

// DeletePayslip deletes a payslip by ID
func (ctrl *PayslipController) DeletePayslip(ctx *gin.Context) {
	id := ctx.Param("payslipId")

	if err := ctrl.DB.Delete(&models.Payslip{}, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete payslip"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
