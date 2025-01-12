package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
	"github.com/wpcodevo/golang-gorm-postgres/middleware"
)

type FileRouteController struct {
	fileController controllers.FileController
}

func NewFileRouteController(fileController controllers.FileController) FileRouteController {
	return FileRouteController{fileController}
}

func (rc *FileRouteController) FileRoute(rg *gin.RouterGroup) {
	router := rg.Group("files")
	router.Use(middleware.DeserializeUser())

	router.POST("/upload", rc.fileController.UploadFile)
	router.GET("/download/:fileName", rc.fileController.DownloadFile)
	router.DELETE("/delete/:fileName", rc.fileController.DeleteFile)
	router.GET("", rc.fileController.ListFiles)
}
