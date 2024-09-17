package db

import (
	"bookServer/internal/domain"
)

type Repository struct {
	Store map[int]domain.Book
	//storeFuture *sql.DB
}

func (d Repository) SaveBookToDataBase(book domain.Book) domain.Book {

	id := len(d.Store) + 1
	d.Store[id] = book

	book.Id = id
	return book
}

func (d Repository) GetBookFromDatabase(id int) domain.Book {
	return d.Store[id]
}
func (d Repository) GetAllBookFromDatabase() map[int]domain.Book {
	return d.Store
}
func (d Repository) DeleteBookFromDatabase(id int) domain.Book {

	delete(d.Store, id)
	return d.Store[id]
}
func (d Repository) UpDateBookToDataBase(book domain.Book, id int) domain.Book {

	d.Store[id] = book

	return book
}
