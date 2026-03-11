package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main(){
	const PORT = "8081"
	mux := new(http.ServeMux)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	homeHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
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
