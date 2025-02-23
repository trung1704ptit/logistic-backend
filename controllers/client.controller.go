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

type ClientController struct {
	DB *gorm.DB
}

func NewClientController(DB *gorm.DB) ClientController {
	return ClientController{DB}
}

func (cc *ClientController) CreateClient(ctx *gin.Context) {
	var payload models.CreateClientRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	newClient := models.Client{
		ID:        uuid.New(),
		Name:      payload.Name,
		Phone:     payload.Phone,
		Address:   payload.Address,
		Note:      payload.Note,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := cc.DB.Create(&newClient)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newClient})
}

func (cc *ClientController) FindClients(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "200")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)

	offset := (intPage - 1) * intLimit

	var clients []models.Client
	result := cc.DB.Limit(intLimit).Offset(offset).Find(&clients)
	if result.Error != nil {

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": clients})
}

func (cc *ClientController) UpdateClient(ctx *gin.Context) {
	clientID := ctx.Param("clientId")

	var payload models.UpdateClient
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var client models.Client
	result := cc.DB.First(&client, "id = ?", clientID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No client with that ID exists"})
		return
	}

	updateData := models.UpdateClient{
		Name:      payload.Name,
		Phone:     payload.Phone,
		Address:   payload.Address,
		Note:      payload.Note,
		UpdatedAt: time.Now(),
	}

	cc.DB.Model(&client).Updates(updateData)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": client})
}

func (cc *ClientController) DeleteClient(ctx *gin.Context) {
	clientID := ctx.Param("clientId")

	result := cc.DB.Delete(&models.Client{}, "id = ?", clientID)
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No client with that ID exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
