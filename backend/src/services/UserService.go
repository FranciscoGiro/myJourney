package services

import (
	"fmt"
	"log"
	"context"
	"errors"
	"time"
	"github.com/FranciscoGiro/myJourney/backend/src/models"
	"github.com/FranciscoGiro/myJourney/backend/src/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


var(
	userNotFoundError = errors.New("User not found")
	emailAlreadyExistsError = errors.New("Email already exists")
	usernameAlreadyExistsError = errors.New("Username already exists")
	unableToRegisterError = errors.New("Unable to register user")
)

type UserService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(user_id string) (models.User, error)
	GetUser(name, email string) (models.User, error)
	CreateUser(name, email, pass string) error
	CheckUserExists(name, email string) error
}

type userService struct {
	userCollection *mongo.Collection
}

func NewUserService() *userService {
	collection := database.GetCollection("Users")
	return &userService{userCollection: collection}
}


func (us *userService) GetUsers(ctx context.Context) ([]models.User, error) {

	result, err := us.userCollection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var users []models.User
	err = result.All(ctx, &users)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

func (us *userService) GetUserByID(user_id string) (models.User, error) {

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uid, _ := primitive.ObjectIDFromHex(user_id)

	err := us.userCollection.FindOne(ctx, bson.M{"_id": uid}).Decode(&user)
	if err != nil {
		fmt.Println("Error:", err)
		return models.User{}, userNotFoundError
	}

	return user, nil
}

func (us *userService) GetUser(name, email string) (models.User, error) {

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{
        "$or",
        bson.A{
            bson.D{
                {"name", name},
            },
            bson.D{
                {"email", email},
            },
        },
    }}


	err := us.userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println("Error:", err)
		return models.User{}, userNotFoundError
	}

	return user, nil
}

func (us *userService) CreateUser(name, email, pass string) error {

	user := &models.User{
		Name: name,
		Email: email,
		Password: pass,
		Role: "BASIC",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	_, err := us.userCollection.InsertOne(ctx, user)
	if err != nil {
		fmt.Println("ERROR:", err)
		return unableToRegisterError
	}

	return nil
}

func (us *userService) CheckUserExists(name, email string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := us.userCollection.CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		fmt.Println("ERROR:", err)
		return err
	}
	if count > 0 {
		return usernameAlreadyExistsError
	}

	count, err = us.userCollection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		fmt.Println("ERROR:", err)
		return err
	}
	if count > 0 {
		return emailAlreadyExistsError
	}

	return nil
}

