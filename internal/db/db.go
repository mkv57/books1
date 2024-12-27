package db

import (
	"books1/internal/domain"
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(rawDB *sql.DB) *Repository {
	return &Repository{db: rawDB}
}

func (d Repository) SaveBookToDataBaseByRAWSql(ctx context.Context, book domain.Book) (*domain.Book, error) {
	book1 := &domain.Book{}
	query := "INSERT INTO books (title, year) VALUES ($1, $2) RETURNING id, title, year"

	err := d.db.QueryRowContext(ctx, query, book.Title, book.Year).Scan(&book1.ID, &book1.Title, &book1.Year)

	if err != nil {
		fmt.Println("error при добавлении книги", err)
	}
	return book1, nil
}

func (d Repository) GetBookFromDatabaseByRAWSql(ctx context.Context, id uint) (*domain.Book, error) {

	book := &domain.Book{}
	query := "SELECT id, title, year FROM books WHERE id = $1"
	err := d.db.QueryRowContext(ctx, query, id).Scan(&book.ID, &book.Title, &book.Year)
	if err != nil {
		fmt.Println("такого ID нет", err)
	}

	return book, nil
}

func (d Repository) GetAllBookFromDatabaseByRAWSql(ctx context.Context) ([]domain.Book, error) {

	books := []domain.Book{}
	query := "SELECT id, title, year FROM books"
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		fmt.Println("error3")
	}
	for rows.Next() {
		book := &domain.Book{}
		rows.Scan(&book.ID, &book.Title, &book.Year)
		books = append(books, *book)
	}

	return books, nil
}

func (d Repository) DeleteBookFromDatabaseByRAWSql(ctx context.Context, id uint) error {
	_, err := d.db.ExecContext(ctx, "DELETE FROM books WHERE id = $1", id)
	if err != nil {
		fmt.Println("ERROR4")
	}
	return nil
}

func (d Repository) UpDateBookToDataBaseByRAWSql(ctx context.Context, book domain.Book) error {
	id := book.ID
	title := book.Title
	year := book.Year
	query := "UPDATE books SET title = $1, year = $2 WHERE id = $3"
	_, err := d.db.ExecContext(ctx, query, title, year, id)
	if err != nil {
		fmt.Println("error5")
	}
	return nil
}
