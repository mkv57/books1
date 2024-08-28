package api

import (
	"bookServer/internal/db"
	"bookServer/internal/domain"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type API struct {
	Di *http.Server `json:"di"`
}

func (p API) GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idstr := r.URL.Query().Get("id")
	idint, err := strconv.Atoi(idstr)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	book, ok := db.Books[idint]
	if !ok {
		handleError(w, http.StatusNotFound, fmt.Errorf("book with id %d not found", idint))
		return
	}
	data, err := json.Marshal(book)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(data)
}

func (p API) AddBook(w http.ResponseWriter, r *http.Request) {
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
	if newBook.Title == "" || newBook.Authors[0] == "" || newBook.Year == 0 {
		//if len(newBook.Title) == 0 || len(newBook.Authors[0]) == 0 || newBook.Year == 0 {    Какой вариант правельный или лучше?

		data, err := json.Marshal("заполнены не все поля")
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		w.Write(data)
		return
	}
	newBook.Id = len(db.Books) + 1 // формируем новый идентификатор
	db.Books[len(db.Books)+1] = newBook
	data, err := json.Marshal(newBook)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(data)
	logger.Info("добавлена книга")
}

func (p API) AllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	limit := query.Get("limit")

	if limit == "" {

		data, err := json.Marshal(db.Books)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}

		w.Write(data)
		logger.Info("отправлен ответ")
		return
	}

	if limit != "" {
		limitNum, err := strconv.Atoi(limit)
		if err != nil {
			handleError(w, http.StatusBadRequest, errors.New("invalid limit parameter"))
			return
		}

		//Проверяем, если параметр limit больше количества книг, то устанавливаем его равным количеству книг
		if limitNum > len(db.Books) {
			limitNum = len(db.Books)
		}

		for i := 1; i <= limitNum; i++ {
			data, err := json.Marshal(db.Books[i])
			if err != nil {
				handleError(w, http.StatusInternalServerError, err)
				return
			}
			w.Write(data)
			logger.Info("отправлен ответ")
		}
		return
	}

}

func (p API) UpdateBook(w http.ResponseWriter, r *http.Request) {
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
	db.Books[book.Id] = book
}

func (p API) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idstr := r.URL.Query().Get("id")
	idint, err := strconv.Atoi(idstr)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	delete(db.Books, idint)
	w.WriteHeader(http.StatusNoContent)
	logger.Info("удалена книга")
}
