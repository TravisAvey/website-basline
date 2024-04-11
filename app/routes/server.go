package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func fileServer(router *mux.Router) {
	fs := http.FileServer(http.Dir("web/dist/"))
	router.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", fs))
}

func Init() {
	router := mux.NewRouter()

	fileServer(router)

	router.HandleFunc("/", index)
	router.HandleFunc("/about", about)
	router.HandleFunc("/blog", blog)
	router.HandleFunc("/blog/posts/new", createPost).Methods("POST")
	router.HandleFunc("/blog/posts/{id}", getPostByID).Methods("GET")
	router.HandleFunc("/blog/posts/s/{slug}", getPostBySlug).Methods("GET")
	router.HandleFunc("/blog/posts/{id}", updatePost).Methods("PUT")
	router.HandleFunc("/blog/posts/{id}", deletePost).Methods("DELETE")
	router.HandleFunc("/blog/posts", getPosts).Methods("GET")
	router.HandleFunc("/gallery", getImages).Methods("GET")
	router.HandleFunc("/gallery/new", newImage).Methods("POST")
	router.HandleFunc("/gallery/{id}", getImage).Methods("GET")
	router.HandleFunc("/gallery/{id}", updateImage).Methods("PUT")
	router.HandleFunc("/gallery/{id}", deleteImage).Methods("DELETE")
	router.HandleFunc("/contact", contact)
	router.HandleFunc("/contact/submit", contactForm).Methods("POST")
	router.HandleFunc("/legal/terms", termsOfUse).Methods("GET")
	router.HandleFunc("/login", login).Methods("GET")
	router.HandleFunc("/dashboard", dashboard)
	router.HandleFunc("/dashboard/posts", dashboardPosts).Methods("GET")
	router.HandleFunc("/dashboard/gallery", dashboardGallery).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)

	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(loggingMiddleware)

	log.Fatal(http.ListenAndServe(":9000", router))
}
