package db

import (
	"books/internal/domain"
	"database/sql"

	//"database/sql"

	"gorm.io/gorm"
)

type Repository struct {
	//Store map[int]domain.Book
	db     *sql.DB
	gormDB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{gormDB: db}
}

func (d Repository) SaveBookToDataBase(book domain.Book) (domain.Book, error) {

	result := d.gormDB.Create(&book)
	if result.Error != nil {
		return domain.Book{}, result.Error
	}
	return book, nil
}

func (d Repository) GetBookFromDatabase(id uint) (domain.Book, error) {
	var book domain.Book
	var result = d.gormDB.First(&book, id)
	//result, _ := d.gormDB.Get("ID =5")
	if result != nil {
		return book, result.Error // ???
	}
	return book, nil
} //
func (d Repository) GetAllBookFromDatabase() ([]domain.Book, error) {
	var books []domain.Book
	result := d.gormDB.Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}

func (d Repository) DeleteBookFromDatabase(id uint) error {
	var book domain.Book
	var result = d.gormDB.Delete(&book, id)
	if result != nil {
		return result.Error
	}
	return nil
}
func (d Repository) UpDateBookToDataBase(book domain.Book) error {
	result := d.gormDB.Save(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
} //
