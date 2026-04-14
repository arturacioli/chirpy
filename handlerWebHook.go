package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) HandleWebHook(w http.ResponseWriter, r *http.Request){
	reqParams := struct{
		Event string `json:"event"`	
		Data struct{
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqParams)
	if err != nil{
		log.Print("Error decoding request")	
		respondWithError(w,http.StatusInternalServerError,"Internal error")
		return
	}
	if reqParams.Event != "user.upgraded"{
		respondWithError(w, http.StatusNoContent, "Wrong event!")
		return
	}
	
	_,err = cfg.db.UpdateUserToRed(r.Context(), reqParams.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found!")
		log.Print(err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")




}
