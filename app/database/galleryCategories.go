package database

// /////////////////////////////////////////////////////////////////
//
//	Gallery Category Crud Ops
//		Gallery Categories: table photo_categories
//		for the categories for the entire gallery
//		not specific to a single image/photo
//
// /////////////////////////////////////////////////////////////////
// NewGalleryCategory Creates new Gallery Category
func NewGalleryCategory(category string) error {
	sqlStatement := `insert into photo_categories(category) values($1);`
	_, err := db.Exec(sqlStatement, category)
	return err
}

func GetGalleryCategories() ([]ImageCategory, error) {
	var cats []ImageCategory
	statement := `select id, category from photo_categories;`
	rows, err := db.Query(statement)
	if err != nil {
		return []ImageCategory{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var id uint64
		var category string

		err = rows.Scan(&id, &category)
		if err != nil {
			return []ImageCategory{}, err
		}
		cat := ImageCategory{
			ID:       id,
			Category: category,
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

func GetGalleryCategory(id uint64) (ImageCategory, error) {
	var category ImageCategory
	statement := `select category from photo_categories where id=$1;`
	rows, err := db.Query(statement, id)
	if err != nil {
		return ImageCategory{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&category.Category)
		if err != nil {
			return ImageCategory{}, err
		}

	}
	category.ID = id
	return category, err
}

func GetGalleryCategoryID(category string) (uint64, error) {
	var id uint64
	statement := `select id from photo_categories where category=$1`
	rows, err := db.Query(statement, category)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

func UpdateGalleryCategory(category ImageCategory) error {
	statement := `update photo_categories set category=$2 where id=$1;`
	_, err := db.Exec(statement, category.ID, category.Category)
	return err
}

func DeleteGalleryCategory(id uint64) error {
	statement := `delete from photo_categories where id=$1;`
	_, err := db.Exec(statement, id)
	return err
}

// /////////////////////////////////////////////////////////////////
//
//	Photo Category Crud Ops
//		Photo Categories: table gallery_categories
//		for the bridge table photo_categories <-> photos
//		for specific categories for an image
//
// /////////////////////////////////////////////////////////////////
// SetPhotoCategory sets a category to a specific image in photo_categories bridge table
func SetPhotoCategory(photoID uint64, categoryID uint64) error {
	statement := `insert into gallery_categories(image_id, category_id) values ($1, (select id from photo_categories where id=$2));`
	_, err := db.Exec(statement, photoID, categoryID)
	return err
}

// GetPhotoCategories returns a list of categories that are linked to the photoID
func GetPhotoCategories(photoID uint64) ([]ImageCategory, error) {
	var categories []ImageCategory
	qwerty := `select photo_categories.id, photo_categories.category 
				from photo_categories, gallery_categories 
				where photo_categories.id = gallery_categories.category_id 
				and gallery_categories.image_id = $1;`
	rows, err := db.Query(qwerty, photoID)
	if err != nil {
		return []ImageCategory{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var categoryID uint64
		var photoCategory string
		err := rows.Scan(&categoryID, &photoCategory)
		if err != nil {
			return []ImageCategory{}, err
		}
		category := ImageCategory{
			ID:       categoryID,
			Category: photoCategory,
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// DeletePhotoCategory deletes the category from the image in the bridge table
func DeletePhotoCategory(photoID uint64, categoryID uint64) error {
	sqlStatement := `delete from gallery_categories where image_id=$1 and category_id=$2;`
	_, err := db.Exec(sqlStatement, photoID, categoryID)
	return err
}
