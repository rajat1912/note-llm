package httpserver

import (
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router chi.Router
}

func New() *Server {
	r := chi.NewRouter()

	r.Get("/auth/{provider}", Provider)
	r.Get("/auth/{provider}/callback", Callback)

	r.Group(func(r chi.Router) {
		r.Use(JWTAuthMiddleware)
		r.Route("/notes", func(r chi.Router) {
			r.Post("/", CreateNoteHandler)
			r.Get("/", GetAllNotesHandler)
			r.Get("/{id}", GetNoteHandler)
			r.Put("/{id}", UpdateNoteHandler)
			r.Delete("/{id}", DeleteNoteHandler)
		})
	})

	return &Server{Router: r}
}
