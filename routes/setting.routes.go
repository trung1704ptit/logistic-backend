package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
	"github.com/wpcodevo/golang-gorm-postgres/middleware"
)

type SettingRouteController struct {
	settingController controllers.SettingController
}

func NewSettingRouteController(settingController controllers.SettingController) SettingRouteController {
	return SettingRouteController{settingController}
}

func (sc *SettingRouteController) SettingRoute(rg *gin.RouterGroup) {
	router := rg.Group("settings")
	router.Use(middleware.DeserializeUser())
	router.POST("", sc.settingController.UpdateSetting)
	router.GET("", sc.settingController.GetSetting)
}
