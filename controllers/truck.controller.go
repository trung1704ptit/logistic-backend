package controllers

import (
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
	return TruckController{DB}
}

// CreateTruck handles creating a new truck record.
func (tc *TruckController) CreateTruck(ctx *gin.Context) {
	var payload *models.CreateTruckRequest

	// Bind the incoming JSON payload to the CreateTruckRequest struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Convert the Contractor string (UUID) to uuid.UUID
	var contractorUUID uuid.UUID
	if payload.Contractor != "" {
		contractorID, err := uuid.Parse(payload.Contractor)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid contractor UUID"})
			return
		}
		contractorUUID = contractorID
	}

	// Set the current time for the truck's CreatedAt and UpdatedAt fields
	now := time.Now()
	newTruck := models.Truck{
		LicensePlate: payload.LicensePlate,
		Capacity:     payload.Capacity,
		Length:       payload.Length,
		Width:        payload.Width,
		Height:       payload.Height,
		Volume:       payload.Volume,
		VehicleType:  payload.VehicleType,
		Brand:        payload.Brand,  // Added brand field
		ContractorID: contractorUUID, // Set the Contractor ID
		Note:         payload.Note,
		Status:       payload.Status, // Set the status
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Create the new truck in the database
	result := tc.DB.Create(&newTruck)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Truck with that license plate already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	// Return the created truck as a response
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newTruck})
}

// UpdateTruck handles updating an existing truck record.
func (tc *TruckController) UpdateTruck(ctx *gin.Context) {
	truckId := ctx.Param("truckId")

	var payload *models.UpdateTruckRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Find the truck by its ID
	var updatedTruck models.Truck
	result := tc.DB.First(&updatedTruck, "id = ?", truckId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Truck not found"})
		return
	}

	// Convert the Contractor string (UUID) to uuid.UUID if provided
	var contractorUUID uuid.UUID
	if payload.Contractor != "" {
		contractorID, err := uuid.Parse(payload.Contractor)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid contractor UUID"})
			return
		}
		contractorUUID = contractorID
	}

	// Set the current time for updating the truck's UpdatedAt field
	now := time.Now()
	truckToUpdate := models.Truck{
		LicensePlate: payload.LicensePlate,
		Capacity:     payload.Capacity,
		Length:       payload.Length,
		Width:        payload.Width,
		Height:       payload.Height,
		Volume:       payload.Volume,
		VehicleType:  payload.VehicleType,
		Brand:        payload.Brand,  // Update brand field
		ContractorID: contractorUUID, // Update contractor UUID
		Note:         payload.Note,
		Status:       payload.Status, // Update the status
		CreatedAt:    updatedTruck.CreatedAt,
		UpdatedAt:    now,
	}

	// Update the truck in the database
	tc.DB.Model(&updatedTruck).Updates(truckToUpdate)

	// Return the updated truck
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedTruck})
}

// FindTruckById handles retrieving a truck by its ID.
func (tc *TruckController) FindTruckById(ctx *gin.Context) {
	truckId := ctx.Param("truckId")

	var truck models.Truck
	result := tc.DB.Preload("Contractor").First(&truck, "id = ?", truckId)
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
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": trucks})
}

// DeleteTruck handles deleting a truck by its ID.
func (tc *TruckController) DeleteTruck(ctx *gin.Context) {
	truckId := ctx.Param("truckId")

	var truck models.Truck
	result := tc.DB.First(&truck, "id = ?", truckId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Truck not found"})
		return
	}

	// Delete the truck from the database
	if err := tc.DB.Delete(&truck).Error; err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
