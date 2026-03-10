package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)
type apiConfig struct {
	fileserverHits atomic.Int32
}
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
func (cfg *apiConfig) handlerMetricsPrinter(w http.ResponseWriter, r *http.Request) {
		header := w.Header()	
		header.Add("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
		`,cfg.fileserverHits.Load())))


}
func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Store(int32(0))
		header := w.Header()	
		header.Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("Count reset to 0"))

}

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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	dat,err := json.Marshal(payload)
	if err != nil{
		log.Printf("Error marshalling ")
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
	w.Write(dat)

}
const PORT = "8081"

func main(){
	mux := new(http.ServeMux)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	homeHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	
	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()	
		header.Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}

	validateChirpHandler := func(w http.ResponseWriter, r *http.Request){
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
			Valid bool `json:"valid"`
		}{
			Valid: true,
		})
	



	}

    mux.Handle("/app/", apiCfg.middlewareMetricsInc(homeHandler))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetricsPrinter)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerMetricsReset)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)



	server := http.Server{
		Addr: ":" + PORT,	
		Handler: mux,
	}

	fmt.Printf("Server running on port: %s\n", PORT)

	err := server.ListenAndServe()
	if err != nil{
		fmt.Print(err)
		os.Exit(1)
	}

}
