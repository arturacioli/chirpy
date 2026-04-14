package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) HandleDeleteChirp(w http.ResponseWriter, r *http.Request){

	pathId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil{
		log.Print("Error parsing uuid")
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	chirp, err := cfg.db.GetSingleChirp(r.Context(), pathId)
	if err != nil{
		log.Print("Error getting chirp")
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	userId := r.Context().Value("userId").(uuid.UUID)
	if chirp.UserID != userId{
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirp.ID)
	if err != nil{
		log.Printf("Error deleting chirp: %v",err)
		respondWithError(w, http.StatusInternalServerError, "Internal server errror")
	}

	w.WriteHeader(204)


}
