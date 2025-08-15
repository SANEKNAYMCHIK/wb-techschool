package httpfiles

import (
	"net/http"

	"github.com/SANEKNAYMCHIK/order-service/internal/cache"
	"github.com/SANEKNAYMCHIK/order-service/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
	cache  *cache.LRUCache
	db     *db.Postgres
}

func NewServer(cache *cache.LRUCache, db *db.Postgres) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	s := &Server{
		router: r,
		cache:  cache,
		db:     db,
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.Get("/order/{order_uid}", s.getOrderHandler)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
