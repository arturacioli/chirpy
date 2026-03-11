package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request){
	type returnVal struct{
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body) 
	params := returnVal{}
	err := decoder.Decode(&params)
	if err != nil{
		log.Printf("Error decoding params: %s\n",err)
		respondWithError(w, 500,"Something went wrong")
	}
	if len(params.Body) > 140 {
		respondWithError(w,400,"Chirp too long")
		return
	}
	if params.Body == "" {
		respondWithError(w, 400, "Body is required")
		return
	}

	respondWithJSON(w,200,struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		CleanedBody: profaneFilter(params.Body),
	})

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

