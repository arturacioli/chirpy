package main

import (
	"fmt"
	"net/http"
)
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

		header := w.Header()	
		header.Add("Content-Type", "text/plain; charset=utf-8")

		if cfg.platform != "dev"{
			w.WriteHeader(403)
			return
		}

		cfg.fileserverHits.Store(int32(0))
		cfg.db.DeleteUsers(r.Context())
		w.WriteHeader(200)
		w.Write([]byte("State reset to 0"))

}
