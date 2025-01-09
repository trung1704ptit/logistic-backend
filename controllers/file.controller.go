package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	BasePath string // Directory where files are stored
}

func NewFileController(BasePath string) FileController {
	return FileController{BasePath}
}

// UploadFile handles file upload
func (fc *FileController) UploadFile(ctx *gin.Context) {
	// Get the file from the request (name it 'file' in the form)
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "could not get the file from the request"})
		return
	}
	defer file.Close()

	if err := os.MkdirAll(fc.BasePath, os.ModePerm); err != nil {
		fmt.Printf("Error creating upload directory: %v\n", err)
	}

	filePath := filepath.Join(fc.BasePath, fileHeader.Filename)

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

// DownloadFile serves a file to the client
func (fc *FileController) DownloadFile(ctx *gin.Context) {
	fileName := ctx.Param("fileName")
	filePath := filepath.Join(fc.BasePath, fileName)

	// Ensure the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "File not found"})
		return
	}

	// Serve the file
	ctx.File(filePath)
}

// DeleteFile removes a file
func (fc *FileController) DeleteFile(ctx *gin.Context) {
	fileName := ctx.Param("fileName")
	filePath := filepath.Join(fc.BasePath, fileName)

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "File not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to delete file"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "File deleted successfully"})
}

// ListFiles returns a list of all files
func (fc *FileController) ListFiles(ctx *gin.Context) {
	files := []string{}

	// Read files from the directory
	err := filepath.Walk(fc.BasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to list files"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "files": files})
}
