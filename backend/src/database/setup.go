package database

import(
	"fmt"
	"os"
	"log"
	"time"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func Init() *mongo.Client {
	// db_username := os.Getenv("MONGODB_USERNAME")
	// db_pass := os.Getenv("MONGODB_PASS")
	// fmt.Println("mongodb+srv://"+db_username+":"+db_pass+"@cluster0.k5ro1ld.mongodb.net/?retryWrites=true&w=majority")

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
    	ApplyURI("mongodb+srv://francisco:123456Kiko@cluster0.k5ro1ld.mongodb.net/?retryWrites=true&w=majority").
    	SetServerAPIOptions(serverAPIOptions)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connect to database")

	return client

}

var client *mongo.Client = Init()

func GetDB() *mongo.Client {
	return client
}


func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client.Disconnect(ctx)
}

func GetCollection(collectionName string) *mongo.Collection{
	db_name := os.Getenv("DB_NAME")
	collection := client.Database(db_name).Collection(collectionName)
	return collection
}