package database

import (
	"fmt"
	"strings"

	"github.com/travisavey/baseline/app/services"
)

func CreateImage(image Image) error {
	photo := image.Image
	categories := image.Categories
	statement := `
	with image as (insert into photos(title, image_url, summary, is_gallery) values ($1, $2, $3, $4) returning id)
	insert into gallery_categories(image_id, category_id) values
	`
	for i, category := range categories {
		statement += "((select id from image), (select id from photo_categories where id=" + fmt.Sprint(category.ID) + "))"
		if i == (len(categories) - 1) {
			statement += ";"
		} else {
			statement += ","
		}
	}
	_, err := db.Exec(statement, photo.Title, photo.ImageURL, photo.Summary, photo.IsGallery)
	return err
}

func GetAllImages() ([]Image, error) {
	var images []Image
	sqlStatement := `select id, image_url, title, summary, is_gallery from photos;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return images, err
	}
	defer rows.Close()
	for rows.Next() {
		var photo Photo
		err := rows.Scan(&photo.ID, &photo.ImageURL, &photo.Title, &photo.Summary, &photo.IsGallery)
		if err != nil {
			return images, err
		}
		categories, catErr := GetPhotoCategories(photo.ID)
		if catErr != nil {
			return images, catErr
		}
		image := Image{
			Image:      photo,
			Categories: categories,
		}
		images = append(images, image)
	}
	return images, nil
}

func GetImage(id uint64) (Image, error) {
	var image Image
	statement := `select image_url, title, summary, is_gallery from photos where id=$1;`
	rows, err := db.Query(statement, id)
	if err != nil {
		return Image{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var photo Photo
		photo.ID = id
		err = rows.Scan(&photo.ImageURL, &photo.Title, &photo.Summary, &photo.IsGallery)
		if err != nil {
			return Image{}, err
		}

		image.Image = photo
		category, err := GetPhotoCategories(id)
		if err != nil {
			return Image{}, err
		}
		image.Categories = category
	}
	return image, nil
}

func UpdateImage(image Image) error {
	photo := image.Image
	updateCategories := image.Categories
	statement := `update photos set title=$2, image_url=$3, summary=$4, is_gallery=$5 where id=$1;`
	_, err := db.Exec(statement, photo.ID, photo.Title, photo.ImageURL, photo.Summary, photo.IsGallery)
	if err != nil {
		return err
	}
	// TODO: update in the future? as of now just removing old categories then adding new/updated categories
	categories, catErr := GetPhotoCategories(photo.ID)
	if catErr != nil {
		return catErr
	}
	for i := range categories {

		delErr := DeletePhotoCategory(photo.ID, categories[i].ID)
		if delErr != nil {
			return delErr
		}
	}
	for i := range updateCategories {

		setErr := SetPhotoCategory(photo.ID, updateCategories[i].ID)
		if setErr != nil {
			return setErr
		}
	}
	return nil
}

func DeleteImage(id uint64) error {
	// remove from S3 storage.
	img, err := GetImage(id)
	if err != nil {
		return err
	}
	url := strings.Split(img.Image.ImageURL, "/")
	key := url[len(url)-1]
	err = services.DeleteImage(key)
	if err != nil {
		return err
	}

	categories, err := GetPhotoCategories(id)
	if err != nil {
		return err
	}

	for _, cat := range categories {
		err = DeletePhotoCategory(id, cat.ID)
		if err != nil {
			return err
		}
	}
	statement := `delete from photos where id=$1;`
	_, err = db.Exec(statement, id)
	return err
}
