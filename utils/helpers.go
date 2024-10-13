package utils

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	    "time"
    "github.com/golang-jwt/jwt/v4"
)


var jwtKey = []byte("c493562f17ef0e93424701ef638ba19572f7b00f4277720c02cec2bc506db9eb9d87633ae573d1826137fb643b952dc24edec18d9bc827234f043fed8a2696cc3bcde605747b4f79e5a08de9027e6e269d15889d3e35fe6fe649b0f976554d4ad7cedaca194783bd6b7a3dfe889a8b4c4cf704cd76265bf5e08017a7bc3d61cc")

// Define the JWT claims
type Claims struct {
    Email string `json:"email"`
    jwt.RegisteredClaims
}


func GenerateJWT(username string) (string, error) {
    // Set expiration time
    expirationTime := time.Now().Add(5 * time.Minute)

    claims := &Claims{
        Email: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign token with secret key
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}


func ValidateJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })    
    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}


func checkValueInMapOfSlice(key string, Data []map[string]string) bool {
	for _, value := range Data {
		if key == value[key] {
			return true
		}
	}
	return false
}

func GetURLParam(r *http.Request, paramName string) (string, error) {

	params := mux.Vars(r)

	param, ok := params[paramName]

	if !ok {
		return "", fmt.Errorf("%s %q", "url parameter not found", paramName)
	}

	return param, nil
}

func GetQueryParam(r *http.Request, paramName string) (string, error) {

	params := r.URL.Query()

	param := params.Get(paramName)

	if param == "" {
		return "", fmt.Errorf("%s %q", "query parameter not found", paramName)
	}

	return param, nil
}
