package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
	"github.com/wpcodevo/golang-gorm-postgres/middleware"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("users")
	router.Use(middleware.DeserializeUser())

	router.GET("/me", uc.userController.GetMe)
	router.GET("", uc.userController.FindUsers)
	router.GET("/:userId", uc.userController.FindUser)
	router.PUT("/:userId", uc.userController.UpdateUser)
	router.DELETE("/:userId", uc.userController.DeleteUser)
}
