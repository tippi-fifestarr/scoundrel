package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tippi-fifestarr/scoundrel/game"
)

// Server represents the API server
type Server struct {
	router         *mux.Router
	handler        *Handler
	sessionManager *game.SessionManager
}

// NewServer creates a new API server
func NewServer() *Server {
	sessionManager := game.NewSessionManager()
	handler := NewHandler(sessionManager)
	router := mux.NewRouter()

	server := &Server{
		router:         router,
		handler:        handler,
		sessionManager: sessionManager,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	// API version prefix
	api := s.router.PathPrefix("/api").Subrouter()

	// Game routes
	api.HandleFunc("/games", s.handler.CreateGameHandler).Methods("POST")
	api.HandleFunc("/games/{id}", s.handler.GetGameHandler).Methods("GET")
	api.HandleFunc("/games/{id}/play/{index}", s.handler.PlayCardHandler).Methods("POST")
	api.HandleFunc("/games/{id}/play-without-weapon/{index}", s.handler.PlayCardWithoutWeaponHandler).Methods("POST")
	api.HandleFunc("/games/{id}/skip", s.handler.SkipRoomHandler).Methods("POST")

	// Root handler
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	// Apply middleware
	s.router.Use(loggingMiddleware)
	s.router.Use(corsMiddleware)
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	srv := &http.Server{
		Handler:      s.router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on %s", addr)
	return srv.ListenAndServe()
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
