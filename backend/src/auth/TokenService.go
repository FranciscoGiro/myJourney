package auth

import(
    "fmt"
    "time"
    "os"
    "errors"
    "github.com/dgrijalva/jwt-go"
    "github.com/FranciscoGiro/myJourney/backend/src/models"
)

var(
    invalidTokenError = errors.New("Invalid authorization token")
    expiredTokenError = errors.New("Authorization token expired")
)


type Payload struct {
	User *models.User `json:user`
	jwt.StandardClaims
}


func GenerateToken(user *models.User) (string, error){

	var secret_key = []byte(os.Getenv("SECRET_KEY"))

    expirationTime := time.Now().Add(15 * time.Minute)

    payload := &Payload{
		User: user,
		StandardClaims: jwt.StandardClaims {
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signed_token, err := token.SignedString(secret_key)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", invalidTokenError
	}

    return signed_token, nil
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
