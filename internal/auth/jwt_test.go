package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJwt(t *testing.T){
	id := uuid.New()
	token,err := MakeJwt(id,"superSecret",time.Hour)	
	if err != nil{
		t.Fatalf("Error making jwt: %v\n",err)
		return
	}

	_,err = ValidateJWT(token, "superSecret")
	if err != nil{
		t.Fatalf("Error validating jwt: %v\n",err)

	}
}

func TestNotValidSecretJwt(t *testing.T){
	id := uuid.New()
	token,err := MakeJwt(id,"superSecret",time.Hour)	
	if err != nil{
		t.Fatalf("Error making jwt: %v\n",err)
		return
	}

	_,err = ValidateJWT(token, "super")
	if err != nil{
		t.Logf("Token rejected for: %v\n",err)
	}


}

func TestExpiredJwt(t *testing.T){
	id := uuid.New()
	token,err := MakeJwt(id,"superSecret",-time.Second)	
	if err != nil{
		t.Fatalf("Error making jwt: %v\n",err)
		return
	}

	_,err = ValidateJWT(token, "superSecret")
	if err != nil{
		t.Logf("Token rejected for: %v\n",err)
	}


}
