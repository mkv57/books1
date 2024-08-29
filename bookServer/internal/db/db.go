package db

import (
	"bookServer/internal/domain"
)

type Repository struct {
	Store map[int]domain.Book
	//storeFuture *sql.DB
}

/*var Books = map[int]domain.Book{
	1: domain.Book{
		Id:      1,
		Title:   "Go на практике",
		Authors: []string{"Мэтт Батчер", "Мэтт Фарина"},
		Year:    2016,
	},
}
*/

func (d Repository) SaveBookToDataBase(book domain.Book) domain.Book {

	id := len(d.Store) + 1
	d.Store[id] = book

	book.Id = id
	return book
}

func (d Repository) GetBookFromDatabase(id int) domain.Book {
	return d.Store[id]
}
