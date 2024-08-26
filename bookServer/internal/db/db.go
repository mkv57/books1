package db

import (
	"bookServer/internal/domain"
)

var Books = map[int]domain.Book{
	1: domain.Book{
		Id:      1,
		Title:   "Go на практике",
		Authors: []string{"Мэтт Батчер", "Мэтт Фарина"},
		Year:    2016,
	},
}

type Server struct {
}
