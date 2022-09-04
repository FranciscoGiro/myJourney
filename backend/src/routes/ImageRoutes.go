package routes


import (
	"github.com/gin-gonic/gin"
	"github.com/FranciscoGiro/myJourney/backend/src/controllers"
)


func ImageRoutes(route *gin.RouterGroup) {

	imageController := controllers.NewImageController()

	route.POST("/images", imageController.UploadImage)

}