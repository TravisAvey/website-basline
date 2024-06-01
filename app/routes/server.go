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
	router.HandleFunc("/blog/posts/s/{slug}", getPostBySlug).Methods("GET")
	router.HandleFunc("/blog/posts", getPosts).Methods("GET")
	router.HandleFunc("/blog/categories", getPosts).Methods("GET")
	router.HandleFunc("/gallery", getImages).Methods("GET")
	router.HandleFunc("/gallery/{id}", getImage).Methods("GET")
	router.HandleFunc("/gallery/{id}", deleteImage).Methods("DELETE")
	router.HandleFunc("/contact", contact)
	router.HandleFunc("/contact/submit", contactForm).Methods("POST")
	router.HandleFunc("/legal/terms", termsOfUse).Methods("GET")
	router.HandleFunc("/login", loginPage).Methods("GET")
	router.HandleFunc("/login", loginAttempt).Methods("POST")
	router.HandleFunc("/logout", logOut).Methods("GET")

	router.HandleFunc("/dashboard", authMiddleware(dashboard))
	router.HandleFunc("/dashboard/posts", authMiddleware(dashboardPosts)).Methods("GET")
	router.HandleFunc("/dashboard/posts/new", createPost).Methods("POST")
	router.HandleFunc("/dashboard/posts/{id}", authMiddleware(getPostByID)).Methods("GET")
	router.HandleFunc("/dashboard/posts/{id}", authMiddleware(deletePost)).Methods("DELETE")
	router.HandleFunc("/dashboard/post/edit/{id}", authMiddleware(editPostView)).Methods("GET")
	router.HandleFunc("/dashboard/post/edit/{id}", authMiddleware(updatePost)).Methods("PUT")
	router.HandleFunc("/dashboard/post/create", authMiddleware(newPost)).Methods("GET")
	router.HandleFunc("/dashboard/blog/count", authMiddleware(dashboardPostCount)).Methods("GET")
	router.HandleFunc("/dashboard/gallery", authMiddleware(dashboardGallery)).Methods("GET")
	router.HandleFunc("/dashboard/gallery/create", authMiddleware(createImageView)).Methods("GET")
	router.HandleFunc("/dashboard/gallery/new", authMiddleware(newImage)).Methods("POST")
	router.HandleFunc("/dashboard/gallery/{id}", authMiddleware(updateImageView)).Methods("GET")
	router.HandleFunc("/dashboard/gallery/{id}", authMiddleware(updateImage)).Methods("PUT")

	router.HandleFunc("/dashboard/messages", authMiddleware(getMessages)).Methods("GET")
	router.HandleFunc("/dashboard/message/{id}", authMiddleware(getMessage)).Methods("GET")
	router.HandleFunc("/dashboard/message/{id}", authMiddleware(messageRead)).Methods("PUT")
	router.HandleFunc("/dashboard/message/{id}", authMiddleware(messageDelete)).Methods("DELETE")
	router.HandleFunc("/dashboard/messages/unread", authMiddleware(getMessageCount)).Methods("GET")
	router.HandleFunc("/sse-messages", authMiddleware(sseEndpoint))

	router.HandleFunc("/sse-login", sseLogin)

	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)

	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(loggingMiddleware)

	log.Fatal(http.ListenAndServe(":9000", router))
}
