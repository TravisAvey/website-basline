package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/travisavey/baseline/app/database"
)

func gallery(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Text string
	}{
		Text: "Gallery Page",
	}

	t, _ := template.ParseFiles("web/templates/pages/gallery.html")
	err := t.Execute(w, data)
	if err != nil {
		w.Write([]byte("Error processing templates.."))
	}
}

// upload/create a image
func newImage(w http.ResponseWriter, r *http.Request) {
	image, err := parseImageData(r)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	err = database.CreateImage(image)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

// get a single image
func getImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	image, err := database.GetImage(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	img := fmt.Sprint(image)
	w.Write([]byte(img))
}

// get all images
func getImages(w http.ResponseWriter, r *http.Request) {
	images, err := database.GetAllImages()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	img := fmt.Sprint(images)
	w.Write([]byte(img))
}

// update a image

// delete a image
