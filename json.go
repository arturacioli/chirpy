package main

import (
	"encoding/json"
	"log"
	"net/http"
)
func respondWithError(w http.ResponseWriter, code int, msg string){
	type returnVal struct{
		Error string `json:"error"`
	}
	error := returnVal{
		Error: msg,
	}
	dat,err := json.Marshal(error)
	if err != nil{
		log.Printf("Error marshalling ")
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any){
	dat,err := json.Marshal(payload)
	if err != nil{
		log.Printf("Error marshalling ")
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
	w.Write(dat)

}	
