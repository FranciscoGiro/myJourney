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

var Client *mongo.Client

func Init() {
	DBusername := os.Getenv("MONGODB_USERNAME")
	DBpass := os.Getenv("MONGODB_PASS")

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
    	ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.k5ro1ld.mongodb.net/?retryWrites=true&w=majority", DBusername, DBpass)).
    	SetServerAPIOptions(serverAPIOptions)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connect to database")
}


func GetDB() *mongo.Client {
	return Client
}


func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	Client.Disconnect(ctx)
}

func GetCollection(collectionName string) *mongo.Collection{
	db_name := os.Getenv("DB_NAME")
	collection := Client.Database(db_name).Collection(collectionName)
	return collection
}