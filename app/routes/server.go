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
	router.HandleFunc("/blog/posts/{id}", getPost).Methods("GET")
	router.HandleFunc("/blog/posts/{id}", updatePost).Methods("PUT")
	router.HandleFunc("/blog/posts/{id}", deletePost).Methods("DELETE")
	router.HandleFunc("/blog/posts", getPosts).Methods("GET")
	router.HandleFunc("/gallery", gallery)
	router.HandleFunc("/contact", contact)
	router.HandleFunc("/dashboard", dashboard)

	router.Use(mux.CORSMethodMiddleware(router))

	log.Fatal(http.ListenAndServe(":9000", router))
}
