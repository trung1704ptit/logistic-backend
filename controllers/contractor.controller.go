package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/models"
	"gorm.io/gorm"
)

type ContractorController struct {
	DB *gorm.DB
}

func NewContractorController(DB *gorm.DB) ContractorController {
	return ContractorController{DB}
}

func (cc *ContractorController) CreateContractor(ctx *gin.Context) {
	var payload *models.CreateContractorRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	newContractor := models.Contractor{
		Name:      payload.Name,
		Phone:     payload.Phone,
		Address:   payload.Address,
		Note:      payload.Note,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if payload.Name == "T&T" || payload.Name == "T & T" {
		newContractor.Type = "internal"
	}

	result := cc.DB.Create(&newContractor)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Contractor with that name already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newContractor})
}

func (cc *ContractorController) UpdateContractor(ctx *gin.Context) {
	contractorId := ctx.Param("contractorId")

	var payload *models.UpdateContractor
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var updatedContractor models.Contractor
	result := cc.DB.Preload("Drivers").Preload("Trucks").First(&updatedContractor, "id = ?", contractorId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Contractor with that ID does not exist"})
		return
	}

	now := time.Now()
	contractorToUpdate := models.Contractor{
		Name:      payload.Name,
		Phone:     payload.Phone,
		Address:   payload.Address,
		Note:      payload.Note,
		CreatedAt: updatedContractor.CreatedAt,
		UpdatedAt: now,
	}

	cc.DB.Model(&updatedContractor).Updates(contractorToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedContractor})
}

func (cc *ContractorController) FindContractorById(ctx *gin.Context) {
	contractorId := ctx.Param("contractorId")

	var contractor models.Contractor
	// Preload drivers and trucks
	result := cc.DB.Preload("Drivers").Preload("Trucks").First(&contractor, "id = ?", contractorId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Contractor with that ID does not exist"})
		return
	}

	// Respond with contractor data including drivers and trucks
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": contractor})
}

func (cc *ContractorController) FindContractors(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "200")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var contractors []models.Contractor
	results := cc.DB.Preload("Drivers").Preload("Trucks").Limit(intLimit).Offset(offset).Order("updated_at DESC").Find(&contractors)
	if results.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": results.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(contractors), "data": contractors})
}

func (cc *ContractorController) DeleteContractor(ctx *gin.Context) {
	contractorId := ctx.Param("contractorId")

	result := cc.DB.Delete(&models.Contractor{}, "id = ?", contractorId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Contractor with that ID does not exist"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
