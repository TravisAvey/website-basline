package routes

import (
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
		// TODO: log error
		w.Write([]byte("Error processing templates.."))
	}
}

// upload/create a image
func newImage(w http.ResponseWriter, r *http.Request) {
	image, err := parseImageData(r)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: Log error
		w.Write([]byte(err.Error()))
	}

	err = database.CreateImage(image)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end--Unable to create the image",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: Log error
		w.Write([]byte(err.Error()))
	}
}

// get a single image
func getImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: log error
		w.Write([]byte(err.Error()))
		return
	}

	image, err := database.GetImage(id)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: Log error
		w.Write([]byte(err.Error()))
		return
	}
	t, _ := template.ParseFiles("web/templates/pages/gallery/image.html")
	err = t.Execute(w, image)
	if err != nil {
		// TODO: log error
		w.Write([]byte(err.Error()))
	}
}

// get all images
func getImages(w http.ResponseWriter, _ *http.Request) {
	images, err := database.GetAllImages()
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: log error
		w.Write([]byte(err.Error()))
		return
	}

	data := struct {
		ImageURL string
		Title    string
		Summary  string
		Images   []database.Image
	}{
		ImageURL: images[0].Image.ImageURL,
		Title:    images[0].Image.Title,
		Summary:  images[0].Image.Summary,
		Images:   images,
	}

	files := getBaseTemplates()
	files = append(files, "web/templates/pages/gallery.html")
	t, _ := template.ParseFiles(files...)
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		// TODO: log error
		w.Write([]byte(err.Error()))
	}
}

// update a image
func updateImage(w http.ResponseWriter, r *http.Request) {
	image, err := parseImageData(r)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: log error
		w.Write([]byte(err.Error()))
		return
	}

	err = database.UpdateImage(image)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: log error
		w.Write([]byte(err.Error()))
	}
}

// delete a image
func deleteImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: log error
		w.Write([]byte(err.Error()))
		return
	}

	err = database.DeleteImage(id)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end--unable to delete image",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: log error
		w.Write([]byte(err.Error()))
	}
}
