package main

import(
	"os"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/FranciscoGiro/myJourney/backend/src/middlewares"
	"github.com/FranciscoGiro/myJourney/backend/src/controllers"
	"github.com/FranciscoGiro/myJourney/backend/src/database"
)

type Handlers struct {
	imageController *controllers.ImageController
	userController *controllers.UserController
}


func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading environment variables")
		os.Exit(1)
	}
	port := os.Getenv("PORT")


	database.Init()
	defer database.Disconnect()

	h := handlers() 
	app := setServer(h)
	app.Run(":"+port)
}



func handlers() *Handlers {
	return &Handlers{
		imageController: controllers.NewImageController(),
		userController: controllers.NewUserController(),
	}
}

func setServer(h *Handlers) *gin.Engine {
	app := gin.Default()
	app.Use(middlewares.Cors())
	
	routes := app.Group("/api/auth")
	routes.POST("/signup", h.userController.Signup)
	routes.POST("/login", h.userController.Login)

	authRoutes := app.Group("/api").Use(middlewares.AuthMiddleware())
	authRoutes.POST("/images", h.imageController.UploadImage)
	authRoutes.GET("/images", h.imageController.GetAllImages)

	authRoutes.GET("/users", h.userController.GetUsers)
	authRoutes.GET("/users/:id", h.userController.GetUser)

	return app

}