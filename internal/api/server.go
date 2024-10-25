package api

import (
	"books/internal/db"
	"books/internal/domain"
	"books/internal/logger"
	"encoding/json"
	"fmt"

	"errors"
	"io"
	"net/http"
	"strconv"
)

type Server struct {
	Database *db.Repository `json:"database"`
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

	//book, err := p.Database.GetBookFromDatabase(uint(idint))
	//if err != nil {
	//	handleError(w, http.StatusInternalServerError, err)
	//}
	//fmt.Println(book)

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

	data, err := json.Marshal(result)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(data)
	log.Info("сохраняем книгу")
} //

func (p Server) AllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	query := r.URL.Query()
	limit := query.Get("limit")

	//var books []domain.Book
	/*Books, err := p.Database.GetAllBookFromDatabase()
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	*/
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

	/*if limitNum > len(*Books) {
			limitNum = len(*Books)
		}

		books = Books[:limitNum]
	} else {
		books = Books
	}
	*/
	data, err := json.Marshal(books)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(data)
	//ctx := r.Context()
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
