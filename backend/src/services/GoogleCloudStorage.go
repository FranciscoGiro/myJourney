package services

import (
	"fmt"
	"context"
	"log"
	"os"
	"io/ioutil"
	"io"
	"time"
	"google.golang.org/api/iterator"
	"strings"
	"errors"
	"cloud.google.com/go/storage"

	"github.com/FranciscoGiro/myJourney/backend/src/database"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	Client = Init()
	UnableToUpload = errors.New("Unable to upload image")
)

func Init() *storage.Client{
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("Unable to connect to Google Cloud Storage. Error", err)
	}

	return client
}


func Close() {
	if err := Client.Close(); err != nil {
		log.Fatal("Unable to close Google Cloud Storage:", err)
	}
}


func GetGCS() *storage.Client {
	return Client
}


//TODO
//should this function continue despite some error?
func UploadImages() error {

	var (
		bucket_name = os.Getenv("BUCKET_NAME")
	)

	files, err := ioutil.ReadDir("src/tmp/")
	if err != nil {
		//erro
		return err
	}


	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()

	for _, f := range files {

		fmt.Println("FILENAME:", f.Name())

		file, err := os.Open(fmt.Sprintf("src/tmp/%s", f.Name()))
		if err != nil {
			fmt.Println("Unable to open file. Error:", err)
			return UnableToUpload
		}


		obj := Client.Bucket(bucket_name).Object("photos/"+file.Name()).NewWriter(ctx)
		if _, err := io.Copy(obj, file); err != nil {
			fmt.Println("Error coping image image to Google Cloud Storage. Error:", err)
			return UnableToUpload
		}

		if err := obj.Close(); err != nil {
			fmt.Println("Error closing Google Cloud Storage file. Error:", err)
			return UnableToUpload
		}

		//isUploaded
		collection := database.GetCollection("Images")

		c, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		i := strings.Split(file.Name(), "-")[1]
		id := strings.Split(i, ".")[0]
		fmt.Println("ID:", id)
		imageID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			fmt.Println("Unable to convert string to ObjectID. Error:", err)
			return UnableToUpload
		}

		fmt.Println("IMAGE ID:", imageID)

		filter := bson.D{{"_id", imageID}}
		update := bson.D{{"$set", bson.D{{"isUploaded", true}}}}

		_, err = collection.UpdateOne(c, filter, update)
		if err != nil {
			fmt.Println("unable to update 'isUploaded'. Error:", err)
			return UnableToUpload
		}

	}

	return nil
}


func GetSignedURLs(userID string) map[string]string {

	var (
		bucketName = os.Getenv("BUCKET_NAME")
		signedURLs = make(map[string]string)
	)

	
	//get all objects with userID prefix

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	prefix := fmt.Sprintf("photos/%s-", userID)
	fmt.Println("Prefix: ", prefix)
	it := Client.Bucket(bucketName).Objects(ctx, &storage.Query{
        Prefix: prefix,
	})

	var objects []string

	for {
        attrs, err := it.Next()
        if err == iterator.Done {
			fmt.Println(err)
                break
        }
        if err != nil {
                return signedURLs //err
        }
		objects = append(objects, attrs.Name)
        fmt.Println("OBJECT:", attrs.Name)
		}

	//sign each object and add to signedURLs

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}

	for _, obj := range objects {
		u, err := Client.Bucket(bucketName).SignedURL(obj, opts)
		if err != nil {
			fmt.Println("SIGNED URLs:", signedURLs)
			//raise error
		}
		imageID := FilenameToImageID(obj)
		signedURLs[imageID] = u
	}

	fmt.Println("SIGNED URLs:", signedURLs)
	return signedURLs
}

func FilenameToImageID(filename string) string{
	i := strings.Split(filename, "-")[1]
	imageID :=  strings.Split(i, ".")[0]
	return imageID
}

