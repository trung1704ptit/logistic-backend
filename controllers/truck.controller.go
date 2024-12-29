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

type TruckController struct {
	DB *gorm.DB
}

func NewTruckController(DB *gorm.DB) TruckController {
	return TruckController{DB}
}

func (tc *TruckController) CreateTruck(ctx *gin.Context) {
	var payload *models.CreateTruckRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Convert contractor string to UUID if present
	var contractorUUID *uuid.UUID
	if payload.Contractor != "" {
		contractorID, err := uuid.Parse(payload.Contractor)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid contractor UUID"})
			return
		}
		contractorUUID = &contractorID
	}

	now := time.Now()
	newTruck := models.Truck{
		LicensePlate: payload.LicensePlate,
		Capacity:     payload.Capacity,
		Length:       payload.Length,
		Width:        payload.Width,
		Height:       payload.Height,
		Volume:       payload.Volume,
		VehicleType:  payload.VehicleType,
		Contractor:   *contractorUUID,
		Note:         payload.Note,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result := tc.DB.Create(&newTruck)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Truck with that license plate already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newTruck})
}

func (tc *TruckController) UpdateTruck(ctx *gin.Context) {
	truckId := ctx.Param("truckId")

	var payload *models.UpdateTruckRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var updatedTruck models.Truck
	result := tc.DB.First(&updatedTruck, "id = ?", truckId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No truck with that ID exists"})
		return
	}

	now := time.Now()
	truckToUpdate := models.Truck{
		LicensePlate: payload.LicensePlate,
		Capacity:     payload.Capacity,
		Length:       payload.Length,
		Width:        payload.Width,
		Height:       payload.Height,
		Volume:       payload.Volume,
		VehicleType:  payload.VehicleType,
		Contractor:   updatedTruck.Contractor,
		Note:         payload.Note,
		CreatedAt:    updatedTruck.CreatedAt,
		UpdatedAt:    now,
	}

	tc.DB.Model(&updatedTruck).Updates(truckToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedTruck})
}

func (tc *TruckController) FindTruckById(ctx *gin.Context) {
	truckId := ctx.Param("truckId")

	var truck models.Truck
	result := tc.DB.First(&truck, "id = ?", truckId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No truck with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": truck})
}

func (tc *TruckController) FindTrucks(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var trucks []models.Truck
	results := tc.DB.Limit(intLimit).Offset(offset).Find(&trucks)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(trucks), "data": trucks})
}

func (tc *TruckController) DeleteTruck(ctx *gin.Context) {
	truckId := ctx.Param("truckId")

	result := tc.DB.Delete(&models.Truck{}, "id = ?", truckId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No truck with that ID exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
