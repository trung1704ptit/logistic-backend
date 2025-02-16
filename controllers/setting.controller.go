package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-gorm-postgres/models"
	"gorm.io/gorm"
)

type SettingController struct {
	DB *gorm.DB
}

func NewSettingController(DB *gorm.DB) SettingController {
	return SettingController{DB: DB}
}

func (sc *SettingController) UpdateSetting(ctx *gin.Context) {
	var payload *models.Setting

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var setting models.Setting
	now := time.Now()

	// Check if a setting exists
	if err := sc.DB.First(&setting).Error; err != nil {
		// If not found, create a new setting
		newSetting := models.Setting{
			ID:        uuid.New(),
			Settings:  payload.Settings,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := sc.DB.Create(&newSetting).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newSetting})
		return
	}

	// If found, update the existing setting
	setting.Settings = payload.Settings
	setting.UpdatedAt = now

	if err := sc.DB.Save(&setting).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": setting})
}

func (sc *SettingController) GetSetting(ctx *gin.Context) {
	var setting models.Setting

	result := sc.DB.Limit(1).Find(&setting)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": setting})
}
