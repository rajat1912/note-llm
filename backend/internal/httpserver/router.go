package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	Router chi.Router
}

func New() *Server {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

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
