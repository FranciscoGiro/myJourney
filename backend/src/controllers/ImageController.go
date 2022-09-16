package controllers

import (
	"fmt"
	"context"
	"time"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/FranciscoGiro/myJourney/backend/src/services"
	"github.com/FranciscoGiro/myJourney/backend/src/models"

	"github.com/gin-gonic/gin"
)

type ImageController struct {
	imageService services.ImageService
}


func NewImageController() *ImageController {
    return &ImageController{imageService: services.NewImageService()}
}


func (ic *ImageController) UploadImage(c *gin.Context){

	user, _ := c.MustGet("user").(models.User)
	userID := user.ID.Hex()
	
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

	image_id, err := ic.imageService.CreateImage(&user,&lat, &lng, &country, &city, &date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}


	err = ic.imageService.StoreImage(&image, user.ID, image_id, file_extension)
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
		err := ic.imageService.UploadImage(&image, userID, image_id, file_extension)
		if err != nil{
			fmt.Println("ERRO A DAR UPLOAD PARA A GCS:", err)
		}
		elapsed := time.Since(start)
    	fmt.Println("File upload took %s", elapsed)
	}()

	c.JSON(http.StatusOK, "ok")

}


func (ic *ImageController) GetAllImages(c *gin.Context){

	user, _ := c.MustGet("user").(models.User)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	images, err := ic.imageService.GetAllImages(ctx, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, images)
}