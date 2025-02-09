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
		FileName:     payload.FileName,
		PriceDetails: payload.Prices,
		OwnerID:      payload.OwnerID,
		OwnerType:    payload.OwnerType,
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

func (pc *PricingController) FindPricingListByContractorID(ctx *gin.Context) {
	ownerId := ctx.Param("ownerId")

	// Use a slice to hold multiple pricings
	var pricings []models.Pricing
	query := pc.DB.Where("owner_id = ?", ownerId).Order("created_at DESC")

	if err := query.Find(&pricings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No pricings found for this contractor"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Return the list of pricings
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": pricings})
}

func (pc *PricingController) FindLatestPricingByContractorID(c *gin.Context) {
	// Parse contractor ID from the request
	ownerId := c.Param("ownerId")
	var latestPricing models.Pricing

	// Query the database to get the latest pricing by contractor ID
	err := pc.DB.Preload("PriceDetails").
		Where("owner_id = ?", ownerId).
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

func (pc *PricingController) FindPricingByContractorIDAndPriceID(c *gin.Context) {
	// Parse contractor ID from the request
	ownerId := c.Param("ownerId")
	priceID := c.Param("priceId")

	var latestPricing models.Pricing

	// Query the database to get the latest pricing by contractor ID
	err := pc.DB.Preload("PriceDetails").
		Where("owner_id = ? and id = ?", ownerId, priceID).
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
func (pc *PricingController) DeleteAllPricingByContractorID(ctx *gin.Context) {
	ownerId := ctx.Param("ownerId")

	tx := pc.DB.Begin()

	if err := tx.Where("owner_id = ?", ownerId).Delete(&models.PriceDetail{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := tx.Where("owner_id = ?", ownerId).Delete(&models.Pricing{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Pricings deleted successfully"})
}

// DeletePricingWithDetails deletes a pricing and all its associated price details by owner_id and pricing_id
func (pc *PricingController) DeletePricingWithDetails(ctx *gin.Context) {
	// Extract owner_id and pricing_id from request parameters
	ownerId := ctx.Param("ownerId")
	priceIDParam := ctx.Param("priceId")

	// Validate UUIDs
	contractorID, err := uuid.Parse(ownerId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid owner_id"})
		return
	}

	pricingID, err := uuid.Parse(priceIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid pricing_id"})
		return
	}

	tx := pc.DB.Begin()

	// Check if the Pricing exists with the given owner_id and pricing_id
	var pricing models.Pricing
	err = tx.Where("id = ? AND owner_id = ?", pricingID, contractorID).First(&pricing).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Pricing not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Delete all associated PriceDetails
	if err := tx.Where("pricing_id = ?", pricingID).Delete(&models.PriceDetail{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete price details"})
		return
	}

	// Delete the Pricing itself
	if err := tx.Delete(&pricing).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete pricing"})
		return
	}

	tx.Commit()
	ctx.JSON(http.StatusNoContent, gin.H{"status": "success", "message": "Pricing and its price details deleted successfully"})
}
