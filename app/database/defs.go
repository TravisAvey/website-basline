package database

import (
	"html/template"

	"github.com/lib/pq"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	ID       uint64 `json:"id"`
}

// /////////////////////////////////////////////////////////////////
//
//	Blog Definitions
//
// /////////////////////////////////////////////////////////////////
// Article the struct to hold a blog post
type Article struct {
	Title       string        `json:"title"`
	ImageURL    string        `json:"imageURL"`
	Summary     string        `json:"summary"`
	Keywords    string        `json:"keywords"`
	Content     string        `json:"content"`
	HTML        template.HTML `json:"html"`
	Slug        string        `json:"slug"`
	DatePosted  pq.NullTime   `json:"datePosted"`
	DateUpdated pq.NullTime   `json:"dateUpdated"`
	ID          int64         `json:"id"`
}

// Category the struct that holds the category type
type Category struct {
	Category string `json:"category"`
	ID       int64  `json:"id"`
}

// Post the article and categories
type Post struct {
	Article    Article    `json:"article"`
	Categories []Category `json:"categories"`
}

// /////////////////////////////////////////////////////////////////
//
//	Gallery Definitions
//
// /////////////////////////////////////////////////////////////////
// Photo the struct that holds the data for the single image
type Photo struct {
	Title     string `json:"title"`
	ImageURL  string `json:"image_url"`
	Summary   string `json:"description"`
	IsGallery bool   `json:"is_gallery"`
	ID        uint64 `json:"id"`
}

type ImageCategory struct {
	Category string `json:"category"`
	ID       uint64 `json:"id"`
}

type Image struct {
	Categories []ImageCategory `json:"categories"`
	Image      Photo           `json:"image"`
}

// /////////////////////////////////////////////////////////////////
//
//	Messages Definitions (for updates for user new subs, etc)
//
// /////////////////////////////////////////////////////////////////
// Message : the object to hold the message
type Message struct {
	Type    string      `json:"type"`
	Header  string      `json:"header"`
	Message string      `json:"message"`
	Sent    pq.NullTime `json:"sentDate"`
	Read    bool        `json:"read"`
	Id      int         `json:"id"`
}
