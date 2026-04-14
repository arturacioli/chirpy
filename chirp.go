package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/arturacioli/chirpy/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig)HandleCreateChirp(w http.ResponseWriter, r *http.Request){
	type reqVal struct{
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body) 
	params := reqVal{}
	err := decoder.Decode(&params)
	if err != nil{
		log.Printf("Error decoding params: %s\n",err)
		respondWithError(w, 500,"Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp too long")
		return
	}

	if params.Body == "" {
		respondWithError(w, 400, "Body is required")
		return
	}

	createChirpParams := database.CreateChirpParams{
		Body: profaneFilter(params.Body),
		UserID: r.Context().Value("userId").(uuid.UUID),
	}
	chirp, err := cfg.db.CreateChirp(r.Context(),createChirpParams)
	if err != nil{
		log.Printf("Error creating chirp %s\n", err)
		respondWithError(w, http.StatusInternalServerError,"Error creating chirp")
	}


	respondWithJSON(w,http.StatusCreated,Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	},
	)


}	


func profaneFilter(s string) string{	
	blackList := []string{"KERFUFFLE","SHARBERT","FORNAX"} 
	words := strings.Fields(s)
	cleaned_body_sl := []string{}
	for _,value := range words{
		if slices.Contains(blackList,strings.ToUpper(value)) {
			value = "****"
		}

		cleaned_body_sl = append(cleaned_body_sl, value)
	}
	return strings.Join(cleaned_body_sl," ")
}

func (cfg *apiConfig)HandleGetChirps(w http.ResponseWriter, r *http.Request){

	if r.URL.Query().Get("author_id") != ""{


	}

	chirps, err := cfg.db.GetChirps(r.Context())

	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
	}
	chirpsSlice := []Chirp{}

	for _,chirp := range(chirps){
		chirpsSlice = append(chirpsSlice, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpsSlice)

}	


func (cfg *apiConfig)HandleGetSingleChirp(w http.ResponseWriter, r *http.Request){
	uuid, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error parsing id")
		return
	}

	chirp, err := cfg.db.GetSingleChirp(r.Context(),uuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error getting chirp")
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})

}



