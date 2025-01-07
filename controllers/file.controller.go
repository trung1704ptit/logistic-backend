package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/initializers"
)

// FileUpload handles the file upload logic and stores the file.
func FileUpload(ctx *gin.Context) {
	config, _ := initializers.LoadConfig(".")

	// Get the file from the request (name it 'file' in the form)
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "could not get the file from the request"})
		return
	}
	defer file.Close()

	// Define the full path where the file will be saved
	// rootPath, _ := os.Getwd()
	// uploadDir := filepath.Join(rootPath, config.UploadFilePath)

	if err := os.MkdirAll(config.UploadFilePath, os.ModePerm); err != nil {
		fmt.Printf("Error creating upload directory: %v\n", err)
	}

	filePath := filepath.Join(config.UploadFilePath, fileHeader.Filename)

	// Create the file at the specified location
	out, err := os.Create(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": fmt.Sprintf("could not create file: %v", err)})
		return
	}
	defer out.Close()

	// Copy the contents of the uploaded file into the newly created file
	_, err = out.ReadFrom(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": fmt.Sprintf("could not save the file: %v", err)})
		return
	}

	// Respond with the file path (or any other necessary information)
	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"filePath": filePath, // You can return a public URL instead of the file path if needed
	})
}
