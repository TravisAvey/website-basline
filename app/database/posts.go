package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

func GetPostCount() (uint64, error) {
	statement := `select count(*) from posts;`
	var count uint64
	row, err := db.Query(statement)
	if err != nil {
		return count, err
	}
	for row.Next() {
		err = row.Scan(&count)
		if err != nil {
			return count, err
		}
	}
	return count, nil
}

func GetAllPosts() ([]Post, error) {
	var posts []Post
	statement := `select id, title, imageurl, summary, keywords, content, slug, dateposted, dateupdated from posts;`
	rows, err := db.Query(statement)
	if err != nil {
		return []Post{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var title, imageUrl, summary, keywords, content, slug string
		var articleID int64
		var datePosted, dateUpdated pq.NullTime
		err = rows.Scan(&articleID, &title, &imageUrl, &summary, &keywords, &content, &slug, &datePosted, &dateUpdated)
		if err != nil {
			return []Post{}, err
		}

		article := Article{
			ID:          articleID,
			Title:       title,
			ImageURL:    imageUrl,
			Summary:     summary,
			Keywords:    keywords,
			Content:     content,
			Slug:        slug,
			DatePosted:  datePosted,
			DateUpdated: dateUpdated,
		}
		categories, caterr := GetPostCategories(articleID)
		if caterr != nil {
			return []Post{}, caterr
		}
		post := Post{
			Article:    article,
			Categories: categories,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostByID(id int64) (Post, error) {
	var article Article
	statement := `select title, summary, slug, imageurl, keywords, dateupdated, dateposted, content from posts where id = $1;`
	article.ID = id
	row, err := db.Query(statement, id)
	if err != nil {
		return Post{}, err
	}
	defer row.Close()
	count := 0
	for row.Next() {
		count++
		err := row.Scan(&article.Title, &article.Summary, &article.Slug, &article.ImageURL, &article.Keywords, &article.DateUpdated, &article.DatePosted, &article.Content)
		if err != nil {
			return Post{}, err
		}
	}
	if count == 0 {
		return Post{}, errors.New("post does not exist")
	}
	categories, err := GetPostCategories(id)
	if err != nil {
		return Post{}, err
	}

	post := Post{
		Article:    article,
		Categories: categories,
	}
	return post, nil
}

func GetPostBySlug(slug string) (Post, error) {
	statement := `select id, title, summary, imageurl, keywords, dateposted, dateupdated, content from posts where slug = $1`
	var article Article
	article.Slug = slug
	row, err := db.Query(statement, slug)
	if err != nil {
		return Post{}, err
	}
	defer row.Close()

	for row.Next() {
		err := row.Scan(&article.ID, &article.Title, &article.Summary, &article.ImageURL, &article.Keywords, &article.DatePosted, &article.DateUpdated, &article.Content)
		if err != nil {
			return Post{}, err
		}
	}
	categories, err := GetPostCategories(article.ID)
	if err != nil {
		return Post{}, err
	}

	post := Post{
		Article:    article,
		Categories: categories,
	}
	return post, nil
}

func GetLastPostInserted() (int, error) {
	var id int
	statement := `select id from posts order by dateposted desc limit 1;`
	rows, err := db.Query(statement)
	if err != nil {
		return id, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func NewPost(post Post) error {
	article := post.Article
	category := post.Categories
	postStatement := `
	with article as (insert into posts (title, dateposted, imageurl, content, summary, keywords, slug) values ($1, $2, $3, $4, $5, $6, $7) returning id) 
	insert into post_categories(post_id, category_id) values
	`
	for i := 0; i < len(category); i++ {
		postStatement += "((select id from article), (select id from categories where id=" + fmt.Sprint(category[i].ID) + "))"
		if i == (len(category) - 1) {
			postStatement += ";"
		} else {
			postStatement += ","
		}
	}
	_, err := db.Exec(postStatement, article.Title, time.Now(), article.ImageURL, article.Content, article.Summary, article.Keywords, article.Slug)
	return err
}

func UpdatePost(post Post) error {
	article := post.Article
	updateCategories := post.Categories
	id := article.ID
	statement := `update posts set title=$2, dateupdated=$3, imageurl=$4, content=$5, summary=$6, keywords=$7, slug=$8 where id=$1;`
	_, err := db.Exec(statement, id, article.Title, pq.NullTime{Time: time.Now()}, article.ImageURL, article.Content, article.Summary, article.Keywords, article.Slug)
	if err != nil {
		return err
	}
	// TODO: update in the future? as of now just removing old categories then adding new/updated categories
	categories, catErr := GetPostCategories(id)
	if catErr != nil {
		return catErr
	}
	for i := range categories {

		delErr := DeletePostCategory(id, categories[i].ID)
		if delErr != nil {
			return delErr
		}
	}
	for i := range updateCategories {

		setErr := SetPostCategory(id, updateCategories[i].ID)
		if setErr != nil {
			return setErr
		}
	}
	return nil
}

func DeletePost(id int64) error {
	// first need to get post_categories
	categories, catErr := GetPostCategories(id)
	if catErr != nil {
		return catErr
	}
	for i := 0; i < len(categories); i++ {
		delErr := DeletePostCategory(id, categories[i].ID)
		if delErr != nil {
			return delErr
		}
	}
	statement := `delete from posts where id = $1;`
	_, err := db.Exec(statement, id)

	return err
}
