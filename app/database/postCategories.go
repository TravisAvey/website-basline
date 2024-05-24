package database

// /////////////////////////////////////////////////////////////////
//
//	Blog Category Crud Ops
//		Blog Categories: table categories
//
// /////////////////////////////////////////////////////////////////
func NewBlogCategory(category Category) error {
	sqlStatement := `insert into categories(category) values($1);`
	_, err := db.Exec(sqlStatement, category.Category)
	return err
}

func GetBlogCategories() ([]Category, error) {
	var categories []Category
	sqlStatement := `select id, category from categories;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return categories, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var category string
		err := rows.Scan(&id, &category)
		if err != nil {
			return categories, err
		}
		blogCategory := Category{
			ID:       id,
			Category: category,
		}
		categories = append(categories, blogCategory)
	}

	return categories, nil
}

func GetBlogCategory(id int64) (Category, error) {
	statement := `select category from categories where id=$1`
	rows, err := db.Query(statement, id)
	if err != nil {
		return Category{}, err
	}
	var category Category
	for rows.Next() {

		err := rows.Scan(&category.Category)
		if err != nil {
			return Category{}, err
		}
	}
	category.ID = id
	return category, nil
}

// EditBlogCategory updates the category with the given category
func EditBlogCategory(category *Category) error {
	sqlStatement := `update categories set category=$2 where id=$1;`
	_, err := db.Exec(sqlStatement, category.ID, category.Category)
	return err
}

// DeleteBlogCategory deletes the category from the db with the id given
func DeleteBlogCategory(id int64) error {
	// first need to remove any references in the post_category table
	postIDs, catErr := getPostCategory(id)
	if catErr != nil {
		return catErr
	}
	if len(postIDs) > 0 {
		for i := range postIDs {
			delErr := DeletePostCategory(postIDs[i], id)
			if delErr != nil {
				return delErr
			}
		}
	}
	// with all references from the post_category bridge table, safe to delete from category
	sqlStatement := `delete from categories where id=$1`
	_, err := db.Exec(sqlStatement, id)
	return err
}

// /////////////////////////////////////////////////////////////////
//
//	Post Category Crud Ops
//		Post Categories: table post_categories
//
// /////////////////////////////////////////////////////////////////

// SetPostCategory creates a new reference for the post and category in the post_category bridge table
func SetPostCategory(postID int64, categoryID int64) error {
	sqlStatement := `insert into post_categories(post_id, category_id) values ($1, (select id from categories where id=$2) );`
	_, err := db.Exec(sqlStatement, postID, categoryID)
	return err
}

// GetPostCategories returns an array of all the post categories from the category table with the post ID given
func GetPostCategories(postID int64) ([]Category, error) {
	var categories []Category
	sqlStatement := `select categories.id, categories.category from categories, post_categories where categories.id = post_categories.category_id and post_categories.post_id = $1;`
	rows, err := db.Query(sqlStatement, postID)
	if err != nil {
		return categories, err
	}
	defer rows.Close()
	for rows.Next() {
		var categoryID int64
		var postCategory string
		err := rows.Scan(&categoryID, &postCategory)
		if err != nil {
			return categories, err
		}
		category := Category{
			ID:       categoryID,
			Category: postCategory,
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func GetPostCategoryID(categoryName string) (int64, error) {
	var id int64
	statement := `select id from categories where category=$1`
	rows, err := db.Query(statement, categoryName)
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

func UpdatePostCategories() error {
	// update categories set category=$2 where id=$1;
	// statement := `update post_categories set post_id=$1 and category_id=$2;`
	return nil
}

// DeletePostCategory deletes the post_category bridge table row with the post ID and category ID given
func DeletePostCategory(postID int64, categoryID int64) error {
	sqlStatement := `delete from post_categories where post_id=$1 and category_id=$2;`
	_, err := db.Exec(sqlStatement, postID, categoryID)
	return err
}

// getPostCategory returns a list of all the post IDs with the given category ID
// This is a helper method when deleting categories when still being referenced in the post_category table
func getPostCategory(categoryID int64) ([]int64, error) {
	var postIDs []int64
	sqlStatement := `select post_id from post_categories where category_id=$1;`
	row, err := db.Query(sqlStatement, categoryID)
	if err != nil {
		return postIDs, err
	}

	for row.Next() {
		var id int64
		err := row.Scan(&id)
		if err != nil {
			return postIDs, err
		}
		postIDs = append(postIDs, id)
	}
	return postIDs, nil
}
