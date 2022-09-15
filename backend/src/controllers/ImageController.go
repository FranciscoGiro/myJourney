package controllers

import (
	"fmt"
	"context"
	"time"
	"errors"
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

	//TODO get real user
	var user_id = "124324235235235235235"
	
	image, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println("Unable to get files from form-data. Error:", err)
		c.JSON(http.StatusInternalServerError, errors.New("Unable to read uploaded files. Please try again"))
	}


	var file_extension = filepath.Ext(header.Filename)


	lat, lng, date, err := ic.imageService.GetMetadata(&image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	//rewind file pointer to beginning
	_, err = image.Seek(0, 0) 
    if err != nil {
		fmt.Println("Unable to rewind file pointer. Error:", err)
		c.JSON(http.StatusInternalServerError, errors.New("Something went wrong. Try again"))
		return
	}

	city, country, err := services.ReverseGeocode(&lat, &lng)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}


	image_id, err := ic.imageService.CreateImage(&lat, &lng, &country, &city, &date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}


	
	err = ic.imageService.StoreImage(&image, user_id, image_id, file_extension)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "ok")

}


func (ic *ImageController) GetAllImages(c *gin.Context){
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := ic.imageService.GetAllImages(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}