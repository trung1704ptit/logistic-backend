package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
	"github.com/wpcodevo/golang-gorm-postgres/initializers"
	"github.com/wpcodevo/golang-gorm-postgres/routes"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	PostController      controllers.PostController
	PostRouteController routes.PostRouteController

	ContractorController      controllers.ContractorController
	ContractorRouteController routes.ContractorRouteController

	DriverController      controllers.DriverController
	DriverRouteController routes.DriverRouteController

	TruckController      controllers.TruckController
	TruckRouteController routes.TruckRouteController

	PricingController      controllers.PricingController
	PricingRouteController routes.PricingRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	PostController = controllers.NewPostController(initializers.DB)
	PostRouteController = routes.NewRoutePostController(PostController)

	ContractorController = controllers.NewContractorController(initializers.DB)
	ContractorRouteController = routes.NewContractorRouteController(ContractorController)

	DriverController = controllers.NewDriverController(initializers.DB)
	DriverRouteController = routes.NewDriverRouteController(DriverController)

	TruckController = controllers.NewTruckController(initializers.DB)
	TruckRouteController = routes.NewTruckRouteController(TruckController)

	PricingController = controllers.NewPricingController(initializers.DB)
	PricingRouteController = routes.NewPricingRouteController(PricingController)

	// Initialize Gin server
	server = gin.Default()
}

func main() {
	// Load the configuration again for server-related settings
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	// Set up CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:5173", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	// Use CORS middleware
	server.Use(cors.New(corsConfig))

	// Define API group and health check route
	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	router.POST("/upload", controllers.FileUpload)

	// Register routes for various controllers
	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	PostRouteController.PostRoute(router)
	PricingRouteController.PricingRoute(router)

	// Register Contractor routes
	ContractorRouteController.ContractorRoute(router)

	// Register Driver routes
	DriverRouteController.DriverRoute(router)

	// Register Truck routes
	TruckRouteController.TruckRoute(router)

	// Start the server
	log.Fatal(server.Run(":" + config.ServerPort))
}
