package auth

import(
    "fmt"
    "time"
    "os"
    "errors"
    "github.com/dgrijalva/jwt-go"
    "github.com/FranciscoGiro/myJourney/backend/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var(
    invalidTokenError = errors.New("Invalid authorization token")
    expiredTokenError = errors.New("Authorization token expired")
)


type Payload struct {
	UserID 		primitive.ObjectID 		`json:userID`
	Username 	string 					`json:username`
	Role 		string 					`json:role`
	User models.User 					`json:user`
	jwt.StandardClaims
}


func GenerateToken(user *models.User) (string, string, error){

	var secret_key = []byte(os.Getenv("SECRET_KEY"))

    token_exp_time := time.Now().Add(5 * time.Minute)
    refresh_exp_time := time.Now().Add(1 * time.Hour)

    payload := &Payload{
		UserID: (*user).ID,
		Username: (*user).Name,
		Role: (*user).Role,
		StandardClaims: jwt.StandardClaims {
			ExpiresAt: token_exp_time.Unix(),
		},
	}

	refresh_payload := &Payload{
		UserID: (*user).ID,
		Username: (*user).Name,
		Role: (*user).Role,
		StandardClaims: jwt.StandardClaims {
			ExpiresAt: refresh_exp_time.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_payload)

	access_token, err := token.SignedString(secret_key)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", "", invalidTokenError
	}

	refresh_token, err := refresh.SignedString(secret_key)
	if err != nil {
		fmt.Println("Error signing refresh token:", err)
		return "", "", invalidTokenError
	}

    return access_token, refresh_token, nil
}

func ValidateToken(token string) (*Payload, error) {

	var secret_key = []byte(os.Getenv("SECRET_KEY"))

    keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, invalidTokenError
		}
		return []byte(secret_key), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, expiredTokenError) {
			return nil, invalidTokenError
		}
		return nil, invalidTokenError
	}

    payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, invalidTokenError
	}

    return payload, nil
}
