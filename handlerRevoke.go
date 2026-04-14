package main

import (
	"net/http"

	"github.com/arturacioli/chirpy/internal/auth"
)

func (cfg *apiConfig)HandlerRevoke(w http.ResponseWriter, r *http.Request){
	token, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error revoking token")
		return
	}


	w.WriteHeader(204)
}
