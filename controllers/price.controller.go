package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-gorm-postgres/models"
	"gorm.io/gorm"
)

type PricingController struct {
	DB *gorm.DB
}

func NewPricingController(DB *gorm.DB) PricingController {
	return PricingController{DB: DB}
}

func (pc *PricingController) CreatePricing(ctx *gin.Context) {
	var payload []models.Pricing // Expecting an array of Pricing objects

	// Bind the incoming JSON payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Start transaction
	tx := pc.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i, pricing := range payload {
		// Validate ContractorID
		if pricing.ContractorID == uuid.Nil {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Contractor ID is required", "index": i})
			return
		}

		// Validate FileName
		if pricing.FileName == "" {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "FileName is required", "index": i})
			return
		}

		// Set metadata fields
		pricing.ID = uuid.New()
		pricing.CreatedAt = time.Now()
		pricing.UpdatedAt = time.Now()

		// Insert the pricing record
		if err := tx.Create(&pricing).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error(), "index": i})
			return
		}

		// Insert nested PriceDetails
		for j, priceDetail := range pricing.Prices {
			priceDetail.ID = uuid.New()          // Generate a new ID for each PriceDetail
			priceDetail.Notes = pricing.FileName // Associate notes with the FileName if applicable

			if err := tx.Create(&priceDetail).Error; err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error(), "index": j})
				return
			}
		}
	}

	tx.Commit()
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Pricing records created successfully"})
}

func (pc *PricingController) FindAllPricing(ctx *gin.Context) {
	contractorID := ctx.Param("contractorId")

	// Parse ContractorID to UUID
	contractorIDUUID, err := uuid.Parse(contractorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid contractor ID"})
		return
	}

	// Get the list of pricing records for the contractor, ordered by CreatedAt
	var prices []models.Pricing
	result := pc.DB.Where("contractor_id = ?", contractorIDUUID).
		Order("created_at DESC"). // Order by CreatedAt in descending order
		Find(&prices)

	// Check for errors
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	// If no records are found
	if len(prices) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No pricing records found for the contractor"})
		return
	}

	// Return the list of pricing records
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": prices})
}

func (pc *PricingController) GetCurrentPricing(ctx *gin.Context) {
	contractorID := ctx.Param("contractorId")

	// Parse ContractorID to UUID
	contractorIDUUID, err := uuid.Parse(contractorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid contractor ID"})
		return
	}

	// Fetch the most recent pricing for the given contractor based on CreatedAt
	var pricing models.Pricing
	result := pc.DB.Where("contractor_id = ?", contractorIDUUID).
		Order("created_at DESC"). // Order by CreatedAt in descending order (most recent first)
		First(&pricing)           // Fetch the most recent pricing record

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No pricing found for the contractor"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		}
		return
	}

	// Return the most recent pricing data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": pricing})
}

func (pc *PricingController) UpdatePricing(ctx *gin.Context) {
	pricingID := ctx.Param("pricingId")
	var payload models.Pricing

	// Bind the JSON payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Fetch the existing pricing record
	var pricing models.Pricing
	if err := pc.DB.Preload("Prices").First(&pricing, "id = ?", pricingID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Pricing not found"})
		return
	}

	// Start transaction
	tx := pc.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update pricing fields
	pricing.ContractorID = payload.ContractorID
	pricing.FileName = payload.FileName
	pricing.UpdatedAt = time.Now()

	if err := tx.Save(&pricing).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Replace PriceDetails
	if err := tx.Where("pricing_id = ?", pricing.ID).Delete(&models.PriceDetail{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	for _, priceDetail := range payload.Prices {
		priceDetail.ID = uuid.New()          // Generate a new ID for each PriceDetail
		priceDetail.Notes = pricing.FileName // Associate notes with the FileName if applicable

		if err := tx.Create(&priceDetail).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
	}

	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": pricing})
}

func (pc *PricingController) DeletePricing(ctx *gin.Context) {
	contractorIDStr := ctx.Param("contractorId")

	// Parse ContractorID to UUID
	contractorID, err := uuid.Parse(contractorIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid Contractor ID"})
		return
	}

	// Delete all Pricing records for the contractor
	result := pc.DB.Where("contractor_id = ?", contractorID).Delete(&models.Pricing{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	// Respond with success message
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "All pricing records deleted successfully"})
}
