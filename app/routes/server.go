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

func Setup() {
	router := mux.NewRouter()

	fileServer(router)

	router.HandleFunc("/", index)
	router.HandleFunc("/about", about)
	router.HandleFunc("/blog", blog)
	router.HandleFunc("/blog/posts/new", createPost).Methods("POST")
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
	router.HandleFunc("/login", loginPage).Methods("GET")
	router.HandleFunc("/login", loginAttempt).Methods("POST")
	router.HandleFunc("/logout", logOut).Methods("GET")

	router.HandleFunc("/dashboard", authMiddleware(dashboard))
	router.HandleFunc("/dashboard/posts", dashboardPosts).Methods("GET")
	router.HandleFunc("/dashboard/posts/{id}", getPostByID).Methods("GET")
	router.HandleFunc("/dashboard/blog/count", dashboardPostCount).Methods("GET")
	router.HandleFunc("/dashboard/gallery", dashboardGallery).Methods("GET")
	router.HandleFunc("/dashboard/messages", getMessages).Methods("GET")
	router.HandleFunc("/dashboard/message/{id}", getMessage).Methods("GET")
	router.HandleFunc("/dashboard/message/{id}", messageRead).Methods("PUT")
	router.HandleFunc("/dashboard/message/{id}", messageDelete).Methods("DELETE")
	router.HandleFunc("/dashboard/messages/unread", getMessageCount).Methods("GET")
	router.HandleFunc("/sse-messages", sseEndpoint)
	router.HandleFunc("/sse-login", sseLogin)

	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)

	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(loggingMiddleware)

	log.Fatal(http.ListenAndServe(":9000", router))
}
