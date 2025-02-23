package controllers

import (
	"net/http"
	"strconv"
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
	var payload models.CreateDriverRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Convert contractor string to UUID if present
	var contractorID uuid.UUID
	if payload.ContractorID != uuid.Nil {
		contractorID = payload.ContractorID
		// Fetch the contractor using the UUID
		var contractor models.Contractor
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
		ContractorID:  contractorID,
		FixedSalary:   payload.FixedSalary,
		Note:          payload.Note,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	result := dc.DB.Create(&newDriver)
	if result.Error != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": result.Error.Error()})
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
	var limit = ctx.DefaultQuery("limit", "200")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var drivers []models.Driver
	results := dc.DB.Limit(intLimit).Offset(offset).Preload("Contractor").Find(&drivers)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(drivers), "data": drivers})
}

func (dc *DriverController) UpdateDriver(ctx *gin.Context) {
	driverId := ctx.Param("driverId")

	var payload models.UpdateDriverRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var driverToUpdate models.Driver
	result := dc.DB.First(&driverToUpdate, "id = ?", driverId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No driver with that ID exists"})
		return
	}

	// Only update fields that are present in the payload
	if payload.FullName != "" {
		driverToUpdate.FullName = payload.FullName
	}
	if payload.Phone != "" {
		driverToUpdate.Phone = payload.Phone
	}
	if payload.CCCD != "" {
		driverToUpdate.CCCD = payload.CCCD
	}
	if !payload.IssueDate.IsZero() {
		driverToUpdate.IssueDate = payload.IssueDate
	}
	if !payload.DateOfBirth.IsZero() {
		driverToUpdate.DateOfBirth = payload.DateOfBirth
	}
	if payload.Address != "" {
		driverToUpdate.Address = payload.Address
	}
	if payload.LicenseNumber != "" {
		driverToUpdate.LicenseNumber = payload.LicenseNumber
	}
	if !payload.LicenseExpiry.IsZero() {
		driverToUpdate.LicenseExpiry = payload.LicenseExpiry
	}
	if payload.ContractorID != uuid.Nil {
		driverToUpdate.ContractorID = payload.ContractorID
	}
	if payload.FixedSalary != nil {
		driverToUpdate.FixedSalary = payload.FixedSalary
	}
	if payload.Note != "" {
		driverToUpdate.Note = payload.Note
	}

	// Update the driver in the database
	result = dc.DB.Model(&driverToUpdate).Updates(driverToUpdate)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": driverToUpdate})
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
