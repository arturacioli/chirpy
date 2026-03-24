package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/arturacioli/chirpy/internal/auth"
)

func (cfg *apiConfig)HandlerLogin(w http.ResponseWriter, r *http.Request){
	reqBody := struct{
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil{
		log.Printf("Error decoding request: %v\n", err)
		respondWithError(w, http.StatusInternalServerError,"Internal Error")
		return
	}
	if reqBody.Email == "" || reqBody.Password == ""{
		respondWithError(w, http.StatusBadRequest,"Not enough arguments, needs to pass email and password!")
		return
	}


	user, err := cfg.db.GetUserByEmail(r.Context(), reqBody.Email)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows) {
        	respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
        	return
        }
		log.Println("Error getting user from db")
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	isEqual,err := auth.CheckPasswordHash(reqBody.Password, user.HashedPassword)
	if err != nil || !isEqual{
       	respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
       	return
	}

	expirationTime := time.Hour

	if reqBody.ExpiresInSeconds > 0 && reqBody.ExpiresInSeconds < 3600 {
	    expirationTime = time.Duration(reqBody.ExpiresInSeconds) * time.Second
	}
	token, err := auth.MakeJwt(user.ID, cfg.secret, expirationTime)
	if err != nil{
		log.Printf("Error creating JWT: %v\n",err)
		respondWithError(w, http.StatusInternalServerError, "Internal error")
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	})



}
