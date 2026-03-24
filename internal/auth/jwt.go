package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJwt(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userId.String(),
	})

	signedToken,err := token.SignedString([]byte(tokenSecret))
	if err != nil{
		return "", err
	}

	return signedToken, nil
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){

	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error){
			return []byte(tokenSecret), nil
		},
	)
	if err != nil{
		return uuid.Nil,err
	}
	issuer,err := token.Claims.GetIssuer()
	if err != nil{
		return uuid.Nil,err
	}
	if issuer != "chirpy-access"{
		return uuid.Nil, jwt.ErrTokenInvalidIssuer 
	}
	
	userId,err := token.Claims.GetSubject()
	if err != nil{
		return uuid.Nil,err
	}
   	uuId,err  := uuid.Parse(userId)
	if err != nil{
		return uuid.Nil,err
	}
	return uuId,nil
}
