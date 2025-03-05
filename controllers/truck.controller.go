package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-gorm-postgres/models"
	"gorm.io/gorm"
)

type TruckController struct {
	DB *gorm.DB
}

func NewTruckController(DB *gorm.DB) TruckController {
	return TruckController{DB: DB}
}

// CreateTruck handles creating a new truck record.
func (tc *TruckController) CreateTruck(ctx *gin.Context) {
	var payload models.CreateTruckRequest

	// Bind JSON payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Parse ContractorID
	var contractorUUID uuid.UUID
	if payload.ContractorID != uuid.Nil {
		contractorUUID = payload.ContractorID
	}

	// Initialize the new truck
	newTruck := models.Truck{
		ID:           uuid.New(),
		LicensePlate: payload.LicensePlate,
		Capacity:     payload.Capacity,
		Length:       payload.Length,
		Width:        payload.Width,
		Height:       payload.Height,
		Volume:       payload.Volume,
		Brand:        payload.Brand,
		ContractorID: contractorUUID,
		Note:         payload.Note,
		Status:       strings.ToLower(payload.Status),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Insert the new truck into the database
	result := tc.DB.Create(&newTruck)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Truck with that license plate already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newTruck})
}

// UpdateTruck handles updating an existing truck record.
func (tc *TruckController) UpdateTruck(ctx *gin.Context) {
	truckID := ctx.Param("truckId")
	var payload models.UpdateTruckRequest

	// Bind JSON payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Fetch the truck record
	var truck models.Truck
	result := tc.DB.First(&truck, "id = ?", truckID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Truck not found"})
		return
	}

	// Update fields only if they are provided
	updates := map[string]interface{}{
		"UpdatedAt": time.Now(),
	}
	if payload.LicensePlate != "" {
		updates["LicensePlate"] = payload.LicensePlate
	}
	if payload.Capacity != nil {
		updates["Capacity"] = payload.Capacity
	}
	if payload.Length != nil {
		updates["Length"] = payload.Length
	}
	if payload.Width != nil {
		updates["Width"] = payload.Width
	}
	if payload.Height != nil {
		updates["Height"] = payload.Height
	}
	if payload.Volume != nil {
		updates["Volume"] = payload.Volume
	}
	if payload.Brand != "" {
		updates["Brand"] = payload.Brand
	}
	if payload.ContractorID != uuid.Nil {
		updates["ContractorID"] = payload.ContractorID
	}
	if payload.Note != "" {
		updates["Note"] = payload.Note
	}
	if payload.Status != "" {
		updates["Status"] = strings.ToLower(payload.Status)
	}

	// Apply updates
	tc.DB.Model(&truck).Updates(updates)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": truck})
}

// FindTruckById handles retrieving a truck by its ID.
func (tc *TruckController) FindTruckById(ctx *gin.Context) {
	truckID := ctx.Param("truckId")

	var truck models.Truck
	result := tc.DB.Preload("Contractor").First(&truck, "id = ?", truckID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Truck not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": truck})
}

// FindTrucks handles retrieving all trucks in the system.
func (tc *TruckController) FindTrucks(ctx *gin.Context) {
	var trucks []models.Truck
	result := tc.DB.Preload("Contractor").Find(&trucks)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": trucks})
}

// DeleteTruck handles deleting a truck by its ID.
func (tc *TruckController) DeleteTruck(ctx *gin.Context) {
	truckID := ctx.Param("truckId")

	if _, err := uuid.Parse(truckID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid truck ID format"})
		return
	}

	var truck models.Truck
	result := tc.DB.First(&truck, "id = ?", truckID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Truck not found"})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	// Perform Soft Delete
	if err := tc.DB.Delete(&truck).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"status": "success", "message": "Truck deleted successfully"})
}
