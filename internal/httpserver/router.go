package httpserver

import (
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router chi.Router
}

func New() *Server {
	r := chi.NewRouter()

	r.Post("/create", CreateNoteHandler)
	return &Server{Router: r}
}
