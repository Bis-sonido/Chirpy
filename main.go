package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Bis-sonido/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform	   string
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	dbPlatform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	const filePath = "."

	cfg := &apiConfig{
		db: dbQueries,
		platform: dbPlatform,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePath)))))

	mux.HandleFunc("GET /admin/metrics", cfg.handlerHits)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	

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
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset endpoint is only available in dev environment"))
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to reset the database: " + err.Error()))
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users table has been reset"))
}
