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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageController struct {
	imageService services.ImageService
}


func NewImageController() *ImageController {
    return &ImageController{imageService: services.NewImageService()}
}


func (ic *ImageController) UploadImage(c *gin.Context){

	userID, _ := c.MustGet("user").(primitive.ObjectID)
	
	image, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println("Unable to get files from form-data. Error:", err)
		c.JSON(http.StatusInternalServerError, errors.New("Unable to read uploaded files. Please try again"))
		return
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

	image_id, err := ic.imageService.CreateImage(&userID, &lat, &lng, &country, &city, &date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}


	err = ic.imageService.StoreImage(&image, userID.Hex(), image_id, file_extension)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	_, err = image.Seek(0, 0) 
    if err != nil {
		fmt.Println("Unable to rewind file pointer. Error:", err)
		c.JSON(http.StatusInternalServerError, errors.New("Something went wrong. Try again"))
		return
	}


	go func() {
		start := time.Now()
		err := ic.imageService.UploadImage(&image, userID.Hex(), image_id, file_extension)
		if err != nil{
			fmt.Println("Error uploading to GCS:", err)
		}
		elapsed := time.Since(start)
    	fmt.Println("File upload took %s", elapsed)
	}()

	c.JSON(http.StatusOK, "ok")

}


func (ic *ImageController) GetAllImages(c *gin.Context){

	userID, _ := c.MustGet("user").(primitive.ObjectID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	images, err := ic.imageService.GetAllImages(ctx, &userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, images)
}