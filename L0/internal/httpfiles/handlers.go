package httpfiles

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *Server) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "order_uid")
	start := time.Now()
	// Finding in cache
	if order, found := s.cache.Get(orderUID); found {
		respondWithJSON(w, http.StatusOK, order)
		log.Printf("Found in cache for %v seconds!", time.Since(start))
		return
	}

	// Check the database
	order, err := s.db.GetOrderByUID(r.Context(), orderUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "order not found")
		log.Printf("Error of find order: %v", err)
		return
	}

	log.Printf("Found in database for %v seconds!", time.Since(start))
	// Updating cache with recently used item
	s.cache.Set(orderUID, order)
	respondWithJSON(w, http.StatusOK, order)
}

// Helper function for answer
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// Helper function for answer with error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
