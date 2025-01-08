package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-gorm-postgres/models"
	"gorm.io/gorm"
)

// PricingController struct
type PricingController struct {
	DB *gorm.DB
}

func NewPricingController(DB *gorm.DB) PricingController {
	return PricingController{DB: DB}
}

// CreatePricing creates a new pricing and its price details
func (pc *PricingController) CreatePricing(ctx *gin.Context) {
	var payload models.CreatePricing

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	pricing := models.Pricing{
		ID:           uuid.New(),
		ContractorID: payload.ContractorID,
		FileName:     payload.FileName,
		PriceDetails: payload.Prices,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tx := pc.DB.Begin()

	if err := tx.Create(&pricing).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	for _, priceDetail := range payload.Prices {
		priceDetail.ID = uuid.New()
		priceDetail.PricingID = pricing.ID

		if err := tx.Create(&priceDetail).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
	}

	tx.Commit()
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": pricing})
}

// GetPricingsByContractorID retrieves pricings and their details for a specific contractor
// GetPricingByContractorID retrieves a single pricing and its details for a specific contractor
func (pc *PricingController) GetPricingByContractorID(ctx *gin.Context) {
	contractorID := ctx.Param("contractorId")

	var pricing models.Pricing
	query := pc.DB.Preload("PriceDetails").Where("contractor_id = ?", contractorID)

	if err := query.First(&pricing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Pricing not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": pricing})
}

func (pc *PricingController) GetLatestPricingByContractorID(c *gin.Context) {
	// Parse contractor ID from the request
	contractorID := c.Param("contractorId")
	var latestPricing models.Pricing

	// Query the database to get the latest pricing by contractor ID
	err := pc.DB.Preload("PriceDetails").
		Where("contractor_id = ?", contractorID).
		Order("created_at DESC").
		First(&latestPricing).Error

	// Handle errors or no record found
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No pricing found for the given contractor ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Return the latest pricing
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": latestPricing})
}

// DeletePricingByContractorID deletes pricings and their price details for a specific contractor
func (pc *PricingController) DeletePricingByContractorID(ctx *gin.Context) {
	contractorID := ctx.Param("contractorId")

	tx := pc.DB.Begin()

	if err := tx.Where("contractor_id = ?", contractorID).Delete(&models.PriceDetail{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := tx.Where("contractor_id = ?", contractorID).Delete(&models.Pricing{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Pricings deleted successfully"})
}
