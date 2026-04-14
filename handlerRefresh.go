package main

import (
	"net/http"
	"time"

	"github.com/arturacioli/chirpy/internal/auth"
)

func (cfg *apiConfig)HandlerRefresh(w http.ResponseWriter, r *http.Request){
	token, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	rToken,err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), rToken.Token)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	newToken, err := auth.MakeJwt(user.ID,cfg.secret,time.Hour)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{
		Token string `json:"token"`
	}{
		Token: newToken,	
	})


}
