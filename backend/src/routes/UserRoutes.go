package routes


import (
	"github.com/gin-gonic/gin"
	"github.com/FranciscoGiro/myJourney/backend/src/controllers"
)


func UserRoutes(route *gin.RouterGroup) {

	userController := controllers.NewUserController()

	route.POST("/auth/signup", userController.Signup)
	route.POST("/auth/login", userController.Login)

	route.GET("/users", userController.GetUsers)
	route.GET("/users/:id", userController.GetUser)

}