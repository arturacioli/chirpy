package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error){
	authString := headers.Get("Authorization")
	if authString == ""{
		return "", errors.New("Missing authorization")
	}
	const prefix = "ApiKey "
	if !strings.HasPrefix(authString, prefix) {
		return "", errors.New("invalid Authorization format")
	}

	return strings.TrimPrefix(authString, prefix), nil
	
}
