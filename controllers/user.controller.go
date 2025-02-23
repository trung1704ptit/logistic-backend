package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wpcodevo/golang-gorm-postgres/models"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := &models.UserResponse{
		ID:        currentUser.ID,
		Name:      currentUser.Name,
		Email:     currentUser.Email,
		Photo:     currentUser.Photo,
		Role:      currentUser.Role,
		Provider:  currentUser.Provider,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (uc *UserController) FindUser(ctx *gin.Context) {
	userID := ctx.Param("userId")
	var user models.User
	if _, err := uuid.Parse(userID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid user ID format"})
		return
	}

	// Find user by ID
	if err := uc.DB.First(&user, "id = ?", userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
}

func (uc *UserController) FindUsers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "200")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit
	var users []models.User
	results := uc.DB.Limit(intLimit).Offset(offset).Find(&users)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": users})
}

// Update user profile
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("userId") // Get user ID from URL
	var user models.User

	if _, err := uuid.Parse(userID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid user ID format"})
		return
	}

	// Find user by ID
	if err := uc.DB.First(&user, "id = ?", userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		return
	}

	// Bind JSON request body to user struct
	var updatedData models.User
	if err := ctx.ShouldBindJSON(&updatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request data"})
		return
	}

	// Update allowed fields
	user.Name = updatedData.Name

	if err := uc.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "User updated successfully", "data": user})
}

// Soft delete user (mark as inactive)
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("userId") // Get user ID from URL
	if _, err := uuid.Parse(userID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid user ID format"})
		return
	}
	var user models.User

	// Find user by ID
	if err := uc.DB.First(&user, "id = ?", userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		return
	}

	// Soft delete user
	if err := uc.DB.Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"status": "success", "message": "User deleted successfully"})
}
