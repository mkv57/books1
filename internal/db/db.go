package db

import (
	"books1/internal/domain"
	"books1/internal/logger"
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

func (d Repository) SaveBookToDataBaseByRAWSql(ctx context.Context, book domain.Book, token string) (*domain.Book, error) {

	userId := "SELECT user_id FROM session WHERE token = $1"
	var userId1 int
	err1 := d.db.QueryRowContext(ctx, userId, token).Scan(&userId1)
	fmt.Println(userId1)
	if err1 != nil {
		fmt.Println("error при добавлении книги 1", err1)
	}
	book1 := &domain.Book{}
	query := "INSERT INTO books (title, year, user_id) VALUES ($1, $2, $3) RETURNING id, title, year"

	err := d.db.QueryRowContext(ctx, query, book.Title, book.Year, userId1).Scan(&book1.ID, &book1.Title, &book1.Year)

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
	log, found := logger.FromContext(ctx)
	if !found {
		log.Debug("GetAllBookFromDatabaseByRAWSql_log")
	}

	books := []domain.Book{}
	query := "SELECT id, title, year FROM books"
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		log.Debug("error3")
	}
	if rows.Err() != nil {
		log.Debug("GetAllBookFromDatabaseByRAWSql")
	}
	for rows.Next() {
		book := &domain.Book{}
		err = rows.Scan(&book.ID, &book.Title, &book.Year)
		books = append(books, *book)
	}
	if err != nil {
		log.Debug("rows.Scan")
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

func (d Repository) SaveUserToDatabase(ctx context.Context, user domain.User) (domain.User, error) {

	user1 := &domain.User{}
	query := "INSERT INTO users (password, email) VALUES ($1, $2) RETURNING user_id, email, password"
	err := d.db.QueryRowContext(ctx, query, user.Password, user.Email).Scan(&user1.ID, &user1.Email, &user1.Password)
	if err != nil {
		fmt.Println("error при добавлении user", err)
	}

	return *user1, nil
}

func (d Repository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	log, found := logger.FromContext(ctx)
	if !found { //!= true {
		log.Debug("GetUserByEmail_log")
	}

	user := &domain.User{}
	query := "SELECT * FROM users WHERE email = $1"
	err := d.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Password, &user.Email)
	if err != nil {
		log.Debug("такого email нет")
	}
	return *user, nil
}
func (d Repository) SaveSessionTodatabase(ctx context.Context, session domain.Session) error {

	session1 := &domain.Session{}
	query := "INSERT INTO session (user_id, token, ip, user_agent, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING  id_session, user_id, token, ip, user_agent, created_at"
	err := d.db.QueryRowContext(ctx, query, session.UserID, session.Token, session.IP, session.UserAgent, session.CreatedAt).
		Scan(&session1.ID, &session1.UserID, &session1.Token, &session1.IP, &session1.UserAgent, &session1.CreatedAt)
	if err != nil {
		fmt.Println("error при добавлении session", err)
	}
	return nil
}
func (d Repository) GetUserByToken(ctx context.Context, token string) (domain.User, error) {

	user := domain.User{}
	query := "select user_id from session where token = $1"
	err := d.db.QueryRowContext(ctx, query, token).Scan(&user.ID)
	if err != nil {
		fmt.Println("такого token нет", err)
	}

	// join users on session.user_id = users.user_id" +
	// "where session.token = $1"
	return user, nil
}
