package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var router *mux.Router

func Init() {
	router = mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Proto)

		w.Write([]byte("<h1>Hello World</h1>"))
	})

	log.Fatal(http.ListenAndServe(":9000", router))
}
