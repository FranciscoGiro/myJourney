package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/FranciscoGiro/myJourney/backend/src/services"

	"github.com/gin-gonic/gin"
)

type ImageController struct {
	imageService services.ImageService
}


func NewImageController() *ImageController {
    return &ImageController{imageService: services.NewImageService()}
}


func (ic *ImageController) UploadImage(c *gin.Context){

	image, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("No file uploaded")
		c.JSON(http.StatusInternalServerError, err)
	}

	var file_extension string
	file_extension = filepath.Ext(header.Filename)


	lat, lng, date, err := ic.imageService.GetMetadata(&image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	city, country, err := services.ReverseGeocode(&lat, &lng)
	if err != nil {
		fmt.Println("Unable to reverse geo coordinates with error:", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}


	image_id, err := ic.imageService.SaveImage(&lat, &lng, &country, &city, &date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	err = ic.imageService.UploadImage(&image, image_id, file_extension)
	if err != nil {
		fmt.Println("Error uploading image:", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "ok")

}
