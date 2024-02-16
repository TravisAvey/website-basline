package routes

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var router *mux.Router

func Init() {
	router = mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Proto)

		data := struct {
			Text string
		}{
			Text: "<h1>Hello, Templates</h1>",
		}

		files := getBaseTemplates()
		files = append(files, "web/templates/index.html")

		t, _ := template.ParseFiles(files...)
		err := t.ExecuteTemplate(w, "base", data)
		if err != nil {
			fmt.Println("Error executing template: ", err.Error())
		}
	})

	log.Fatal(http.ListenAndServe(":9000", router))
}

// TODO: move this method to another file in the future
// maybe a utils file, or another if there is
// functionality needed for the go templating
func getBaseTemplates() []string {
	return []string{
		"web/templates/base.html",
		"web/templates/partials/nav.html",
		"web/templates/partials/footer.html",
	}
}
