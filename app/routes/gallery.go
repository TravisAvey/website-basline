package routes

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/travisavey/baseline/app/database"
	"github.com/travisavey/baseline/app/services"
)

type imageCats struct {
	Category string
	Selected bool
}

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

func createImageView(w http.ResponseWriter, _ *http.Request) {
	cats, err := database.GetGalleryCategories()
	if err != nil {
		msg := getResponseMsg("Failed to get the images", Error)
		sendSSEMessage(msg)
		// TODO: log error
		return
	}
	data := struct {
		Categories []database.ImageCategory
	}{
		Categories: cats,
	}

	t, _ := template.ParseFiles("web/templates/pages/dashboard/new-image.html")
	err = t.Execute(w, data)
	if err != nil {
		// TODO: Log error
		msg := getResponseMsg("There was an error creating the page", Error)
		sendSSEMessage(msg)
		w.Write([]byte(err.Error()))
	}
}

// upload/create a image
func newImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(64 << 20)
	_, header, err := r.FormFile("image")
	if err != nil {
		// TODO: log error
		// TODO: send sendSSEMessage
		fmt.Println("r.FormFile err:", err.Error())
		return
	}
	fmt.Println(header.Filename)

	isGallery, err := strconv.ParseBool(r.FormValue("is-gallery"))
	if err != nil {
		// TODO: log error
		// TODO: send sendSSEMessage
		fmt.Println("ParseBool err:", err.Error())
		return
	}

	categories, err := parseImageCategories(r)
	if err != nil {
		// TODO: log error
		// TODO: send sendSSEMessage
		fmt.Println("parseImageCategories err:", err.Error())
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		// TODO: log error
		// TODO: send sse msg
		fmt.Println("FormFile error:", err.Error())
		return
	}

	defer file.Close()
	filename := header.Filename
	tmpImg := fmt.Sprintf("./temp/uncmp-%v", filename)
	output := fmt.Sprintf("./temp/%v", filename)

	out, err := os.Create(tmpImg)
	if err != nil {
		fmt.Println("create tempImg error:", err.Error())
		// TODO: log error
		// TODO: send sse msg
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println("Copy file error:", err.Error())
		// TODO: log error
		// TODO: send sse msg
		return
	}

	err = services.CompressImage(tmpImg, output)
	if err != nil {
		fmt.Println("CompressImage err:", err.Error())
		// TODO: log error
		// TODO: send sse msg
		return
	}

	err = services.SendImage(output, filename)
	if err != nil {
		fmt.Println("SendImage err:", err.Error())
		// TODO: log error
		// TODO: send sse msg
		return
	}

	image := database.Image{
		Image: database.Photo{
			ImageURL:  services.GetS3Url() + "/" + filename,
			Title:     r.FormValue("title"),
			Summary:   r.FormValue("description"),
			IsGallery: isGallery,
		},
		Categories: categories,
	}

	err = database.CreateImage(image)
	if err != nil {
		// TODO: Log error
		// TODO: send sse msg
		fmt.Println("createImage err:", err.Error())
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

// a func just for all post categories and a flag for an active category
func checkImageCategories(allCats, cats []database.ImageCategory) []imageCats {
	var categories []imageCats

	for _, cat := range allCats {
		curCats := imageCats{
			Category: cat.Category,
			Selected: false,
		}
		for _, c := range cats {
			if c.Category == cat.Category {
				curCats.Selected = true
			}
		}
		categories = append(categories, curCats)
	}

	return categories
}

func updateImageView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		// TODO: log error
		// TODO: send sse msg
		fmt.Println("ParseUint error:", err.Error())
		return
	}

	var img database.Image
	img, err = database.GetImage(id)
	if err != nil {
		// TODO: log error
		// TODO: send sse msg
		fmt.Println("GetImage error:", err.Error())
		return
	}

	cats, err := database.GetGalleryCategories()
	if err != nil {
		msg := getResponseMsg("Failed to get the images", Error)
		sendSSEMessage(msg)
		// TODO: log error
		return
	}

	data := struct {
		Categories []imageCats
		Image      database.Photo
	}{
		Image:      img.Image,
		Categories: checkImageCategories(cats, img.Categories),
	}

	t, _ := template.ParseFiles("web/templates/pages/dashboard/edit-image.html")
	err = t.Execute(w, data)
	if err != nil {
		// TODO: log error
		// TODO: send sse msg
		fmt.Println("ParseFiles error:", err.Error())
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
