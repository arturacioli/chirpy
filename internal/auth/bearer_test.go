package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T){
	authToken := "imatoken"
	header := http.Header{} 
	header.Add("Authorization","Bearer " + authToken)

	token, err := GetBearerToken(header)
	if err != nil{
		t.Fatalf("Error %s\n",err)
	}

	if token != authToken{
		t.Fatalf("%s shoud be equal to %s",token, authToken)
	}
}
