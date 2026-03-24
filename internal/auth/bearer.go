package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header)(string,error){
	tokenString := headers.Get("Authorization")
	if tokenString == ""{
		return "", errors.New("Missing authorization")
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(tokenString, prefix) {
		return "", errors.New("invalid Authorization format")
	}

	return strings.TrimPrefix(tokenString, prefix), nil
	
}
