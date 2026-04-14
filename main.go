package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"

	"github.com/arturacioli/chirpy/internal/auth"
	"github.com/arturacioli/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
	secret string
	apiKey string
}


func main(){
	godotenv.Load()
	const PORT = "8081"

	dbURL := os.Getenv("DB_URL")
	if dbURL == ""{
		log.Fatalf("db url env variable missing")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == ""{
		log.Fatalf("secret env variable missing")
	}
	platform := os.Getenv("PLATFORM")
	if platform == ""{
		log.Fatalf("platform env variable missing")
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == ""{
		log.Fatalf("api key env variable missing")
	}

	db, err := sql.Open("postgres", dbURL)

	dbQueries := database.New(db)
	if err != nil {
		fmt.Print(err)	
		os.Exit(1)
	}

	mux := new(http.ServeMux)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		platform: platform, 
		apiKey: apiKey,
	}

	homeHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(homeHandler))

	mux.HandleFunc("GET /api/healthz", healthHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetricsPrinter)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerMetricsReset)

	mux.HandleFunc("POST /api/users", apiCfg.HandleCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.HandlerEditUser)

	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogin)
	
	mux.HandleFunc("POST /api/chirps", apiCfg.middlewareAuth(apiCfg.HandleCreateChirp))
	mux.HandleFunc("GET /api/chirps", apiCfg.HandleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandleGetSingleChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.middlewareAuth(apiCfg.HandleDeleteChirp))

	mux.HandleFunc("POST /api/refresh", apiCfg.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandlerRevoke)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.middlewareApiAuth(apiCfg.HandleWebHook))


	server := http.Server{
		Addr: ":" + PORT,	
		Handler: mux,
	}

	fmt.Printf("Server running on port: %s\n", PORT)

	err = server.ListenAndServe()
	if err != nil{
		fmt.Print(err)
		os.Exit(1)
	}

}

func (cfg *apiConfig) middlewareAuth(next http.HandlerFunc) http.HandlerFunc{

	return func(w http.ResponseWriter, r *http.Request){
		token,err := auth.GetBearerToken(r.Header)
		if err != nil{
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		id, err := auth.ValidateJWT(token, cfg.secret)
		if err != nil{
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		
		ctx := context.WithValue(r.Context(), "userId", id)
		next(w, r.WithContext(ctx))

	}
}

func (cfg *apiConfig) middlewareApiAuth(next http.HandlerFunc) http.HandlerFunc{

	return func(w http.ResponseWriter, r *http.Request){
		token,err := auth.GetApiKey(r.Header)
		if err != nil || token != cfg.apiKey{
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		next(w, r)

	}
}

