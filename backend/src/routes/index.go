package routes


import (
	"github.com/gin-gonic/gin"
)


func AllRoutes(route *gin.RouterGroup) {

	ImageRoutes(route)
	UserRoutes(route)
}