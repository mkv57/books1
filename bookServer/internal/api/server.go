package api

import (
	"books/bookServer/internal/db"
	"books/bookServer/internal/domain"
	"encoding/json"

	"errors"
	"io"
	"net/http"
	"strconv"
)

type Server struct {
	Database *db.Repository `json:"database"`
}

func (p Server) GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idstr := r.URL.Query().Get("id")
	idint, err := strconv.Atoi(idstr)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}

	book, err := p.Database.GetBookFromDatabase(uint(idint))
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
	}

	data, err := json.Marshal(book)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(data)
} //

func (p Server) AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	result, err := p.Database.SaveBookToDataBase(newBook)
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
} //

func (p Server) AllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	limit := query.Get("limit")

	var books []domain.Book
	Books, err := p.Database.GetAllBookFromDatabase()
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
		if limitNum > len(Books) {
			limitNum = len(Books)
		}

		books = Books[:limitNum]
	} else {
		books = Books
	}
	data, err := json.Marshal(books)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(data)
	Logger.Info("отправлен ответ")
} //

func (p Server) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	err = p.Database.UpDateBookToDataBase(book)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}
} //

func (p Server) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idstr := r.URL.Query().Get("id")
	idint, err := strconv.Atoi(idstr)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	err = p.Database.DeleteBookFromDatabase(uint(idint))
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
	}
	w.WriteHeader(http.StatusNoContent)
	Logger.Info("удалена книга")
} //
