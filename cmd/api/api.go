package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/fixztter/chi-api/service/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(r)

	log.Printf("listening on port%s\n", s.addr)

	return http.ListenAndServe(s.addr, r)
}
