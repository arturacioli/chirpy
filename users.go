package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig)HandleCreateUser(w http.ResponseWriter, r *http.Request){
	type reqJson struct {
		Email string `json:"email"`
	}	
	
	decoder := json.NewDecoder(r.Body)
	params := reqJson{}

	err := decoder.Decode(&params)
	if err != nil{
		log.Println("Error decoding request body")
		respondWithError(w,500,"Something went wrong")
		return
	}
	
	
	user, err := cfg.db.CreateUser(r.Context(),params.Email)
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
	})



}
