package db

import (
	"books/internal/domain"
	"context"
	"database/sql"
	"fmt"

	//"database/sql"

	"gorm.io/gorm"
)

type Repository struct {
	//Store map[int]domain.Book
	db     *sql.DB
	gormDB *gorm.DB
}

func NewRepository(db2 *gorm.DB, rawDB *sql.DB) *Repository {
	return &Repository{gormDB: db2, db: rawDB}
}

func (d Repository) SaveBookToDataBase(book domain.Book) (domain.Book, error) {

	result := d.gormDB.Create(&book)
	if result.Error != nil {
		return domain.Book{}, result.Error
	}
	return book, nil
}

func (d Repository) SaveBookToDataBaseByRAWSql(ctx context.Context, book domain.Book) (domain.Book, error) {

	err := d.db.QueryRowContext(ctx, "insert into books (title, year) values ($1, $2)returning *", book.Title, book.Year)
	if err != nil {
		fmt.Println("ERROR")
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
}

func (d Repository) GetBookFromDatabaseByRAWSql(ctx context.Context, id uint) (domain.Book, error) {
	var book domain.Book
	err := d.db.QueryRowContext(ctx, "select id, title, year from books where id = $1, id").Scan(&book.ID, &book.Title, &book.Year)
	if err != nil {
		fmt.Println("ERROR")
	}
	return book, nil
}

func (d Repository) GetAllBookFromDatabase() ([]domain.Book, error) {
	var books []domain.Book
	result := d.gormDB.Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}

func (d Repository) GetAllBookFromDatabaseByRAWSql(ctx context.Context) ([]domain.Book, error) {
	var books []domain.Book
	rows, err := d.db.QueryContext(ctx, "select * from books")
	if err != nil {
		fmt.Println("error")
	}
	for rows.Next() {
		var book domain.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Year)
		books = append(books, book)
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

func (d Repository) DeleteBookFromDatabaseByRAWSql(ctx context.Context, id uint) error {
	_, err := d.db.ExecContext(ctx, "delete from books where id = $1", id)
	if err != nil {
		fmt.Println("ERROR")
	}
	return nil
}

func (d Repository) UpDateBookToDataBase(book domain.Book) error {
	result := d.gormDB.Save(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d Repository) UpDateBookToDataBaseByRAWSql(ctx context.Context, book domain.Book) error {
	err := d.db.QueryRowContext(ctx, "update books set title = $1, year = $2 where id = $3", "ttt", 2024)
	if err != nil {
		fmt.Println("error")
	}
	return nil
}
