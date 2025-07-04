package httpserver

import (
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router chi.Router
}

func New() *Server {
	r := chi.NewRouter()

	r.Post("/notes", CreateNoteHandler)
	r.Get("/notes/{id}", GetNoteHandler)
	return &Server{Router: r}
}
