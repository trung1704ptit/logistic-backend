package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-gorm-postgres/models"
	"gorm.io/gorm"
)

type DriverController struct {
	DB *gorm.DB
}

func NewDriverController(DB *gorm.DB) DriverController {
	return DriverController{DB}
}

func (dc *DriverController) CreateDriver(ctx *gin.Context) {
	var payload *models.CreateDriverRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Convert contractor string to UUID if present
	var contractor models.Contractor
	if payload.Contractor != "" {
		contractorID, err := uuid.Parse(payload.Contractor)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid contractor UUID"})
			return
		}
		// Fetch the contractor using the UUID
		result := dc.DB.First(&contractor, "id = ?", contractorID)
		if result.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Contractor not found"})
			return
		}
	}

	now := time.Now()
	newDriver := models.Driver{
		FullName:      payload.FullName,
		Phone:         payload.Phone,
		CCCD:          payload.CCCD,
		IssueDate:     payload.IssueDate,
		DateOfBirth:   payload.DateOfBirth,
		Address:       payload.Address,
		LicenseNumber: payload.LicenseNumber,
		LicenseExpiry: payload.LicenseExpiry,
		ContractorID:  contractor.ID, // Store only the ContractorID, not the full object
		Note:          payload.Note,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	result := dc.DB.Create(&newDriver)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Driver with that license number already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newDriver})
}

func (dc *DriverController) FindDriverById(ctx *gin.Context) {
	driverId := ctx.Param("driverId")

	var driver models.Driver
	result := dc.DB.Preload("Contractor").First(&driver, "id = ?", driverId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No driver with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": driver})
}

func (dc *DriverController) FindDrivers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var drivers []models.Driver
	results := dc.DB.Limit(intLimit).Offset(offset).Preload("Contractor").Find(&drivers)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(drivers), "data": drivers})
}

func (dc *DriverController) UpdateDriver(ctx *gin.Context) {
	driverId := ctx.Param("driverId")

	var payload *models.UpdateDriverRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var updatedDriver models.Driver
	result := dc.DB.First(&updatedDriver, "id = ?", driverId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No driver with that ID exists"})
		return
	}

	// Convert contractor string to UUID if present
	var contractor models.Contractor
	if payload.Contractor != "" {
		contractorID, err := uuid.Parse(payload.Contractor)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid contractor UUID"})
			return
		}
		// Fetch the contractor using the UUID
		result := dc.DB.First(&contractor, "id = ?", contractorID)
		if result.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Contractor not found"})
			return
		}
	}

	now := time.Now()
	driverToUpdate := models.Driver{
		FullName:      payload.FullName,
		Phone:         payload.Phone,
		CCCD:          payload.CCCD,
		IssueDate:     payload.IssueDate,
		DateOfBirth:   payload.DateOfBirth,
		Address:       payload.Address,
		LicenseNumber: payload.LicenseNumber,
		LicenseExpiry: payload.LicenseExpiry,
		ContractorID:  contractor.ID, // Store only the ContractorID
		Note:          payload.Note,
		CreatedAt:     updatedDriver.CreatedAt,
		UpdatedAt:     now,
	}

	// Perform the update
	dc.DB.Model(&updatedDriver).Updates(driverToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedDriver})
}

func (dc *DriverController) DeleteDriver(ctx *gin.Context) {
	driverId := ctx.Param("driverId")

	result := dc.DB.Delete(&models.Driver{}, "id = ?", driverId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No driver with that ID exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
