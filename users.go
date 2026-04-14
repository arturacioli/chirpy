package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arturacioli/chirpy/internal/auth"
	"github.com/arturacioli/chirpy/internal/database"
)

func (cfg *apiConfig)HandleCreateUser(w http.ResponseWriter, r *http.Request){
	type reqJson struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}	
	
	decoder := json.NewDecoder(r.Body)
	params := reqJson{}

	err := decoder.Decode(&params)
	if err != nil{
		log.Println("Error decoding request body")
		respondWithError(w,500,"Something went wrong")
		return
	}
	
	hashedPassword,err := auth.HashPassword(params.Password)	
	if err != nil{
		log.Printf("Error hashing password: %v",err)
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(),database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil{
		log.Printf("Error creating user: %v", err) 
		respondWithError(w,500,"Something went wrong")
		return
	}
	
	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	})



}
