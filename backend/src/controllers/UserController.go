package controllers

import (
	"fmt"
	"os"
	"context"
	"time"
	"net/http"
	"github.com/FranciscoGiro/myJourney/backend/src/services"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}


func NewUserController() *UserController {
    return &UserController{userService: services.NewUserService()}
}


func (uc *UserController) Signup(c *gin.Context) {

	var body struct {
		Name string	`json:"name"`
		Email string	`json:"email"`
		Password string	`json:"password"`
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}


	err = uc.userService.CheckUserExists(body.Name, body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hash_pass, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10) // SALT ROUNDS
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash given password"})
		return
	}

	err = uc.userService.CreateUser(body.Name, body.Email, string(hash_pass)) //provavel erro na hash
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})

}

func (uc *UserController) Login(c *gin.Context) {
	var body struct {
		Name string	`json:"name"`
		Password string	`json:"password"`
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request body"})
		return
	}

	user, err := uc.userService.GetUser(body.Name, body.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}


	secret_key := os.Getenv("SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
		"name": user.Name,
		"role": user.Role,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	signed_token, err := token.SignedString(secret_key)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing token"})
	}

	c.SetSameSize(http.SameSiteLaxMode)
	c.Cookie("Authorization", signed_token, 3600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (uc *UserController) GetUsers(c *gin.Context){

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := uc.userService.GetUsers(ctx)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)

}

func (uc *UserController) GetUser(c *gin.Context){
	user_id := c.Param("id")
	fmt.Println("User ID:", user_id)

	user, err := uc.userService.GetUserByID(user_id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, user)

}

