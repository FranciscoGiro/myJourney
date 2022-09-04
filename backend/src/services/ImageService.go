package services

import (
	"fmt"
	"time"
	"context"
	"errors"
	"io"
	"os"
	"github.com/FranciscoGiro/myJourney/backend/src/database"
	"github.com/FranciscoGiro/myJourney/backend/src/models"
	"github.com/rwcarlsen/goexif/exif"
	"mime/multipart"


	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageService interface {
	SaveImage(lat *float64, lng *float64, country *string, city *string, date *time.Time) (string, error)
	UploadImage(image *multipart.File, id, file_extension string) error
	GetMetadata(image *multipart.File) (float64, float64, time.Time, error)
}

type imageService struct {
	imageCollection *mongo.Collection
}

func NewImageService() *imageService {
	collection := database.GetCollection("Images")
	return &imageService{imageCollection: collection}
}


func (is *imageService) SaveImage(lat *float64, lng *float64, 
	country *string, city *string, date *time.Time) (string, error) {

	id := primitive.NewObjectID()

	newImage := &models.Image{
		ID: id,
		User_id: "user id", // TODO , user needs to be inserted 
		City: *city,
		Country: *country,
		Lat: *lat,
		Lng: *lng,
		Date: *date,
		CreatedAt: time.Now(), //fix date format
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := is.imageCollection.InsertOne(ctx, newImage)
	if err != nil {
		fmt.Printf("Error inserting image:", err)
		return "", err
	}

	return id.Hex(), nil
} 


func (is *imageService) GetMetadata(image *multipart.File) (float64, float64, time.Time, error){

	var(
		lat, lng float64
		date time.Time
	)

	metadata, err := exif.Decode(*image)
    if err != nil {
		// fix this, should not be error 400
		fmt.Println("Unable to read EXIF metadata: ", err)
		return lat, lng, date, errors.New("Unable to extract image metadata")
    }

	lat, lng, err = metadata.LatLong()
	if err != nil {
		fmt.Println("Unable to retrieve coordinates from image:", err)
		return lat, lng, date, errors.New("No geo coordinates found in image")
	}

	date, err = metadata.DateTime()
	if err != nil {
		fmt.Println("Unable to retrieve date from image:", err)
		return lat, lng, date, errors.New("Unable to retrieve date from image")
	}

	return lat, lng, date, nil
}


func (is *imageService) UploadImage(image *multipart.File, id, file_extension string) error {

	var (
		bucket_name = os.Getenv("BUCKET_NAME")
		client = GetGCS()
		filename = fmt.Sprintf("%s%s",id, file_extension)
	)

	fmt.Println(filename)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	obj := client.Bucket(bucket_name).Object("photos/"+filename).NewWriter(ctx)
	if _, err := io.Copy(obj, *image); err != nil {
		fmt.Println("Error coping image file", err)
		return err
	}
	if err := obj.Close(); err != nil {
		fmt.Println("Error closing image file", err)
		return err
	}

	return nil
}