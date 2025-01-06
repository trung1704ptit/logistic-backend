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

// CreatePricing handles creating a new pricing record.
func (pc *PricingController) CreatePricing(ctx *gin.Context) {
	var payload models.Pricing

	// Bind JSON payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Convert contractorID to uuid.UUID
	contractorID, err := uuid.Parse(ctx.Param("contractorId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid contractor ID"})
		return
	}

	// Set version for new pricing
	if payload.Version == "" {
		payload.Version = time.Now().Format(time.RFC3339) // Set version as current timestamp
	}

	// Set IsCurrent to true for new pricing
	payload.IsCurrent = true
	payload.ContractorID = contractorID // Set the contractor ID after converting it to UUID
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()

	// Insert the new pricing into the database
	result := pc.DB.Create(&payload)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": payload})
}

// FindAllPricing handles retrieving all pricing records for a contractor.
func (pc *PricingController) FindAllPricing(ctx *gin.Context) {
	contractorID := ctx.Param("contractorId")

	var prices []models.Pricing
	result := pc.DB.Where("contractor_id = ?", contractorID).Find(&prices)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	if len(prices) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No pricing found for the contractor"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": prices})
}

// GetCurrentPricing retrieves the current pricing version for a specific contractor.
func (pc *PricingController) GetCurrentPricing(ctx *gin.Context) {
	contractorID := ctx.Param("contractorId")

	var pricing models.Pricing
	result := pc.DB.Where("contractor_id = ? AND is_current = ?", contractorID, true).First(&pricing)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No current pricing found for the contractor"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": pricing})
}

// UpdatePricing handles updating a pricing record by its ID.
func (pc *PricingController) UpdatePricing(ctx *gin.Context) {
	pricingID := ctx.Param("pricingId")
	var payload models.Pricing

	// Bind JSON payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Fetch the pricing record
	var pricing models.Pricing
	result := pc.DB.First(&pricing, "id = ?", pricingID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Pricing not found"})
		return
	}

	// Update fields if provided
	updates := map[string]interface{}{"UpdatedAt": time.Now()}

	if payload.FromLocationCity != "" {
		updates["FromLocationCity"] = payload.FromLocationCity
	}
	if payload.FromLocationDistrict != "" {
		updates["FromLocationDistrict"] = payload.FromLocationDistrict
	}
	if payload.ToLocationCity != "" {
		updates["ToLocationCity"] = payload.ToLocationCity
	}
	if payload.ToLocationDistrict != "" {
		updates["ToLocationDistrict"] = payload.ToLocationDistrict
	}
	if len(payload.Prices) > 0 {
		updates["Prices"] = payload.Prices
	}
	if payload.Note != "" {
		updates["Note"] = payload.Note
	}

	// Apply updates
	pc.DB.Model(&pricing).Updates(updates)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": pricing})
}

// DeletePricing handles deleting all pricing records for a contractor.
func (pc *PricingController) DeletePricing(ctx *gin.Context) {
	contractorIDStr := ctx.Param("contractorId")

	// Convert contractorID to uuid.UUID
	contractorID, err := uuid.Parse(contractorIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid contractor ID"})
		return
	}

	// Delete all pricing records for the contractor
	result := pc.DB.Where("contractor_id = ?", contractorID).Delete(&models.Pricing{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	// Respond with success
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Pricing records deleted successfully"})
}
