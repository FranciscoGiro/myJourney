package services

import (
	"context"
	"log"
	"cloud.google.com/go/storage"
)



func Init() *storage.Client{
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("Unable to connect to Google Cloud Storage:", err)
	}

	return client
}

var Client = Init()


func Close() {
	if err := Client.Close(); err != nil {
		log.Fatal("Unable to close Google Cloud Storage:", err)
	}
}


func GetGCS() *storage.Client {
	return Client
}
