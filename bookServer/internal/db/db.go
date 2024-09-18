package db

import (
	"bookServer/internal/domain"

	"gorm.io/gorm"
)

type Repository struct {
	//Store map[int]domain.Book
	//storeFuture *sql.DB
	gormDB *gorm.DB
}

func NewRepository(db *gorm.Db) *Repository {
	return &Repository{gormDB: db}
}

func (d Repository) SaveBookToDataBase(book domain.Book) (domain.Book, error) {

	result := d.gormDB.Creat(&book)
	if result.Error != nil {
		return domain.Book{}, result.Error
	}
	return book, nil
}

func (d Repository) GetBookFromDatabase(id uint) domain.Book {
	//return d.Store[id]
	return domain.Book{}
}
func (d Repository) GetAllBookFromDatabase() []domain.Book {
	//return d.Store
	return nil
}
func (d Repository) DeleteBookFromDatabase(id uint) domain.Book {

	//delete(d.Store, id)
	//return d.Store[id]
	return domain.Book{}
}
func (d Repository) UpDateBookToDataBase(book domain.Book, id uint) domain.Book {

	//d.Store[id] = book

	return book
}
