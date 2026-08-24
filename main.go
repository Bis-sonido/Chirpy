package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/Bis-sonido/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var badWords = []string{"kerfuffle", "sharbert", "fornax"}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	const filePath = "."

	cfg := &apiConfig{
		db: dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePath)))))

	mux.HandleFunc("GET /admin/metrics", cfg.handlerHits)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerJson)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePath, server.Addr)
	listenErr := server.ListenAndServe()
	if listenErr != nil {
		log.Println("Error starting server:", listenErr)
	}

}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	hit := cfg.fileserverHits.Load()
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
	</html>
	`, hit)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	fmt.Fprintf(w, "Hits reset to 0\n")
}

func handlerJson(w http.ResponseWriter, r *http.Request) {
	type jsonRequest struct {
		Body  string `json:"body"`
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	var params jsonRequest
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Post must be less than 140 characters")
		return
	}

	splitBody := strings.Split(params.Body, " ")
	for i, word := range splitBody {
		for _, badWord := range badWords {
			if strings.ToLower(word) == strings.ToLower(badWord) {
				// Replace the bad word with asterisks
				splitBody[i] = "****"
			}
		}
	}
	params.Body = strings.Join(splitBody, " ")

	type jsonResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	respBody := jsonResponse{
		CleanedBody: params.Body,
	}

	respondWithJSON(w, http.StatusOK, respBody)
}
