package services

import (
	"fmt"
	"time"
	"context"
	"errors"
	"os"
	"io"
	"github.com/FranciscoGiro/myJourney/backend/src/database"
	"github.com/FranciscoGiro/myJourney/backend/src/models"
	"github.com/rwcarlsen/goexif/exif"
	"mime/multipart"


	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var(
	unableToReadEXIF = errors.New("Unable to read image properties. Check if image contains it, otherwise upload with location")
	unableToReadDate = errors.New("Unable to read image original date. Check if image contains it, otherwise upload with location and date")
	unableToUpload = errors.New("Unable to upload image")
	unableToFindImages = errors.New("No images found")
)

type ImageService interface {
	GetAllImages(ctx context.Context, user *models.User) ([]ImageInfo, error)
	CreateImage(user *models.User, lat *float64, lng *float64, country *string, 
		city *string, date *time.Time) (string, error)
	UploadImage(image *multipart.File, userID, imageID, file_extension string) error
	StoreImage(image *multipart.File, user_id primitive.ObjectID, image_id, file_extension string) error
	GetMetadata(image *multipart.File) (float64, float64, time.Time, error)
}

type imageService struct {
	imageCollection *mongo.Collection
}

type ImageInfo struct {
	ID 			string 		`json:"id"`
	Url 		string 		`json:"url"`
	City 		string 	    `json:"city"`
	Country 	string 	    `json:"country"`
	Lat 		float64 	`json:"lat"`
	Lng 		float64 	`json:"lng"`
	Date	    time.Time	`json:"date"`
}

func NewImageService() *imageService {
	collection := database.GetCollection("Images")
	return &imageService{imageCollection: collection}
}

func (is *imageService) GetMetadata(image *multipart.File) (float64, float64, time.Time, error){

	var(
		lat, lng float64
		date time.Time
	)

	metadata, err := exif.Decode(*image)
    if err != nil {
		// fix this, should not be error 400
		fmt.Println("Unable to read EXIF metadata. Error:", err)
		return lat, lng, date, unableToReadEXIF
    }

	lat, lng, err = metadata.LatLong()
	if err != nil {
		fmt.Println("Unable to retrieve LAT LONG from image. Error:", err)
		return lat, lng, date, unableToReadEXIF
	}

	date, err = metadata.DateTime()
	if err != nil {
		fmt.Println("Unable to retrieve DATE from image. Error:", err)
		return lat, lng, date, unableToReadDate
	}

	return lat, lng, date, nil
}

//save image inside tmp folder
//save it also inside a temporary user folder in order to upload to google cloud storage
func (is *imageService) CreateImage(user *models.User, lat *float64, lng *float64, 
	country *string, city *string, date *time.Time) (string, error) {

	//TODO id should be given by MongoDb and then retrieved once uploaded
	id := primitive.NewObjectID()

	fmt.Println("CITY2:", *city)
	fmt.Println("COUNTRY2:", *country)
	
	newImage := &models.Image{
		ID: id,
		User_id: (*user).ID,
		City: *city,
		Country: *country,
		Lat: *lat,
		Lng: *lng,
		Date: *date,
		IsUploaded: false,
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := is.imageCollection.InsertOne(ctx, newImage)
	if err != nil {
		fmt.Printf("Error inserting image in database. Error:", err)
		return "", unableToUpload
	}

	return id.Hex(), nil
} 




//TODO 
//change this function to a scheduled one
//goes through tmp folder and uploads all images in there
func (is *imageService) UploadImage(image *multipart.File, userID, imageID, file_extension string) error {

	var (
		bucket_name = os.Getenv("BUCKET_NAME")
		client = GetGCS()
		filename = fmt.Sprintf("%s-%s%s",userID, imageID, file_extension)
	)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	obj := client.Bucket(bucket_name).Object("photos/"+filename).NewWriter(ctx)
	if _, err := io.Copy(obj, *image); err != nil {
		fmt.Println("Error coping image image to Google Cloud Storage. Error:", err)
		return unableToUpload
	}
	if err := obj.Close(); err != nil {
		fmt.Println("Error closing Google Cloud Storage file. Error:", err)
		return unableToUpload
	}

	//TODO isUploaded= true in database

	return nil
}

func (is *imageService) StoreImage(image *multipart.File, user_id primitive.ObjectID, image_id, file_extension string) error {

	dst, err := os.Create(fmt.Sprintf("src/tmp/%s-%s%s", user_id.Hex(), image_id, file_extension))
	defer dst.Close()
	if err != nil {
		fmt.Println("Unable to create temp file. Error:", err)
		return unableToUpload
	}

	if _, err := io.Copy(dst, *image); err != nil {
		fmt.Println("Unable to copy image to temp file. Error:", err)
		return unableToUpload
	}

	return nil
}

//TODO
//needs to look to database and google cloud storage and create a struct in order to
//return the right information to frontend
func (is *imageService) GetAllImages(ctx context.Context, user *models.User) ([]ImageInfo, error) {


	//confirm that it returns
	fmt.Println("USER_ID:", (*user).ID)
	result, err := is.imageCollection.Find(ctx, bson.M{"user_id": (*user).ID})
	if err != nil {
		fmt.Println("Error retrieving images from database. Error:", err)
		return nil, unableToFindImages
	}

	var images []models.Image
	err = result.All(ctx, &images)
	if err != nil {
		fmt.Println("Error parsing images from database. Error:", err)
		return nil, unableToFindImages
	}

	if images == nil{
		return []ImageInfo{}, nil
	}


	// get all images from user in gcs
	userID := ((*user).ID).Hex()
	signedURLs := GetSignedURLs(userID)

	//loop through images in db and match with signed url

	var res []ImageInfo

	for _, image := range images {
		// TODO
		// check that in fact image id exists in gcs
		imageID := image.ID.Hex()

		imageInfo := ImageInfo{
			ID: imageID,
			Url: signedURLs[imageID],
			City: image.City,
			Country: image.Country,
			Lat: image.Lat,
			Lng: image.Lng,
			Date: image.Date,
		}

		res = append(res, imageInfo)
	}

	// return this new structure with required info

	return res, nil
}
