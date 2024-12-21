package api

import (
	"books1/internal/domain"
	"books1/internal/logger"
	"context"
	"encoding/json"
	"fmt"

	"errors"
	"io"
	"net/http"
	"strconv"
)

type Store interface {
	SaveBookToDataBaseByRAWSql(ctx context.Context, book domain.Book) (*domain.Book, error)
	GetBookFromDatabaseByRAWSql(ctx context.Context, id uint) (*domain.Book, error)
	GetAllBookFromDatabaseByRAWSql(ctx context.Context) ([]domain.Book, error)
	DeleteBookFromDatabaseByRAWSql(ctx context.Context, id uint) error
	UpDateBookToDataBaseByRAWSql(ctx context.Context, book domain.Book) error
}

type Server struct {
	//Database *db.Repository `json:"database"`
	Database Store
}

func (p Server) GetBook(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log, found := logger.FromContext(ctx)
	if found == false {
		handleError(w, http.StatusInternalServerError, errors.New("Проблемы у нас"))
		return
	}
	log.Info("отправляем в ответ книгу")
	w.Header().Set("Content-Type", "application/json")
	idstr := r.URL.Query().Get("id")
	idint, err := strconv.Atoi(idstr)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}

	book1, err := p.Database.GetBookFromDatabaseByRAWSql(ctx, uint(idint))
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
	}

	data, err := json.Marshal(book1)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(data)

}

func (p Server) AddBook(ctx context.Context, request *pb.AddBookRequest) (*pb.AddBookResponse, error) {
	_, found := logger.FromContext(ctx)
	if found == false {
		return nil, errors.New("нет логера")
	}

	newBook := domain.Book{
		Title: request.Title,
		Year:  int(request.Year),
	}

	result, err := p.Database.SaveBookToDataBaseByRAWSql(ctx, newBook)
	if err != nil {
		return nil, err
	}

	return &pb.AddBookResponse{Book: &pb.Book{
		Id:    int64(result.ID),
		Title: result.Title,
		Year:  int32(result.Year),
	}}, nil

	//log.Info("сохраняем книгу")
}

/*
func (p Server) AddBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log, found := logger.FromContext(ctx)
	if found == false {
		handleError(w, http.StatusInternalServerError, errors.New("Проблемы у нас"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	r.Context()
	jsong, err := io.ReadAll(r.Body)

	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	var newBook domain.Book
	err = json.Unmarshal(jsong, &newBook)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}

	result, err := p.Database.SaveBookToDataBaseByRAWSql(ctx, newBook)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	fmt.Println(result)

	data, err := json.Marshal(result)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(data)
	log.Info("сохраняем книгу")
} //
*/

func (p Server) AllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	query := r.URL.Query()
	limit := query.Get("limit")

	books, err := p.Database.GetAllBookFromDatabaseByRAWSql(ctx)
	//fmt.Println(books)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	if limit != "" {
		limitNum, err := strconv.Atoi(limit)
		if err != nil {
			handleError(w, http.StatusBadRequest, errors.New("invalid limit parameter"))
			return
		}
		fmt.Println(limitNum)
	}

	data, err := json.Marshal(books)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(data)

	log, found := logger.FromContext(ctx)
	if found == false {
		handleError(w, http.StatusInternalServerError, errors.New("Проблемы у нас"))
		return
	}
	log.Info("отправлен ответ")
}

func (p Server) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	var book domain.Book
	err = json.Unmarshal(data, &book)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	err = p.Database.UpDateBookToDataBaseByRAWSql(ctx, book)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}

	log, found := logger.FromContext(ctx)
	if found == false {
		handleError(w, http.StatusInternalServerError, errors.New("Проблемы у нас"))
		return
	}
	log.Info("обновили книгу")
}

func (p Server) DeleteBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log, found := logger.FromContext(ctx)
	if found == false {
		handleError(w, http.StatusInternalServerError, errors.New("Проблемы у нас"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	idstr := r.URL.Query().Get("id")
	idint, err := strconv.Atoi(idstr)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	err = p.Database.DeleteBookFromDatabaseByRAWSql(ctx, uint(idint))
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
	}
	w.WriteHeader(http.StatusNoContent)
	log.Info("удалена книга")
}
