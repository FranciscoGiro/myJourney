package controllers

import (
	"fmt"
	"context"
	"time"
	"net/http"
	"errors"
	"github.com/FranciscoGiro/myJourney/backend/src/services"
	"github.com/FranciscoGiro/myJourney/backend/src/models"
	"github.com/FranciscoGiro/myJourney/backend/src/auth"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController() *UserController {
    return &UserController{
		userService: services.NewUserService(),
	}
}

type Claims struct {
	User *models.User `json:user`
	jwt.StandardClaims
}

func (uc *UserController) Signup(c *gin.Context) {

	var body struct {
		Name 		string		`json:"username"`
		Email 		string		`json:"email"`
		Password 	string		`json:"password"`
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
		fmt.Println("Unable to generate hash password. Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash given password"})
		return
	}

	err = uc.userService.CreateUser(body.Name, body.Email, hash_pass) //provavel erro na hash
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})

}

func (uc *UserController) Login(c *gin.Context) {
	var body struct {
		Name string	`json:"username"`
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

	fmt.Println("Hashed:", user.Password)
	fmt.Println("New:", body.Password)
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	access_token, refresh_token, err := auth.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.New("Unable to register. Please try again"))
		return
	}

	uc.userService.SaveRefreshToken(refresh_token, user.ID)

	c.SetCookie("Authorization", refresh_token, 3600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"access_token": access_token})
}

func (uc *UserController) Refresh(c *gin.Context) {

	token, err := c.Cookie("Authorization")
	if err != nil {
		fmt.Println("Unable to retrieve auth cookie. Error:", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("Unable to read auth cookie"))
		return
	}

	payload, err := auth.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("Invalid token"))
			return
		}

	givenUser := (*payload).User
	fmt.Println("UserID:", givenUser.ID.Hex())
	user, err := uc.userService.GetUserByID(givenUser.ID.Hex()) //TODO error Hex import
	if err != nil {
		fmt.Println("O erro está aqui", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if token != user.RefreshToken {
		fmt.Println("error here", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("Invalid token"))
		return
	}

	// everything ok. Now generate new tokens

	access_token, refresh_token, err := auth.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.New("Unable to generate tokens"))
		return
	}

	uc.userService.SaveRefreshToken(refresh_token, user.ID)

	c.SetCookie("Authorization", refresh_token, 3600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"access_token": access_token})
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
	userID := c.Param("id")

	user, err := uc.userService.GetUserByID(userID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, user)

}

