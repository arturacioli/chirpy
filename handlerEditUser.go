package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arturacioli/chirpy/internal/auth"
	"github.com/arturacioli/chirpy/internal/database"
)

func (cfg *apiConfig) HandlerEditUser(w http.ResponseWriter, r *http.Request){
	reqBody := struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}{}


	token, err := auth.GetBearerToken(r.Header)	
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}


	userId,err := auth.ValidateJWT(token, cfg.secret)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized,err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	hashedPassword, err := auth.HashPassword(reqBody.Password)
	if err != nil{
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	updatedUser,err := cfg.db.UpdateUserEmailAndPassword(
		r.Context(),
		database.UpdateUserEmailAndPasswordParams{
			ID: userId,
			Email: reqBody.Email,
			HashedPassword: hashedPassword,
		})
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}


	respondWithJSON(w, http.StatusOK, User{
		ID: updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email: updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed.Bool,
	})
	
}
