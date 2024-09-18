package domain

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title string `json:"title"`
	//Authors []string `json:"authors"`
	Year int `json:"year"`
}
