package router

import (
	"net/http"
	"strings"
	"log"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func New(port string) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/public/v1", func(r chi.Router) {
		// config to all public routes

		// specific routes
		r.Route("/posts", func(r chi.Router) {
			r.Get("/", listPosts)
			// r.Post("/", createPost)

			// r.Route("/{postID}", func(r chi.Router) {
			// 	r.Use(postCtx)
			// 	r.Get("/", getPost)
			// 	r.Put("/", updatePost)
			// 	r.Delete("/", deletePost)
			// })
		})
	})

	// r.Route("/api/admin/v1", func(r chi.Router) {

	// })

	r.Route("/", func(root chi.Router) {
		fileServer(root, "", "/dist/", http.Dir("assets/public/dist/"))
		fileServer(root, "", "/", http.Dir("assets/public/static/"))
	})

	log.Printf("Up and running on port %s...", port)
	http.ListenAndServe(":" + port, r)
}

func fileServer(r chi.Router, basePath string, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("fileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(basePath+path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func listPosts(w http.ResponseWriter, r *http.Request) {}