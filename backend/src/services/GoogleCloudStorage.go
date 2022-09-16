package services

import (
	"fmt"
	"context"
	"log"
	"os"
	"time"
	"google.golang.org/api/iterator"
	"strings"
	"errors"
	"cloud.google.com/go/storage"
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
			
			fmt.Println("Error signing object. Error:", err)
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

