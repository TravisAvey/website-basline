package routes

import (
	"fmt"
	"html/template"
	"net/http"
)

func about(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Text string
	}{
		Text: "About Page",
	}

	t, _ := template.ParseFiles("web/templates/pages/about.html")
	err := t.Execute(w, data)
	if err != nil {
		errStr := fmt.Sprint("error: ", err.Error())
		w.Write([]byte(errStr))
	}
}
