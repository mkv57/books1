package api

import (
	"bookServer/internal/db"
	"bookServer/internal/domain"
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

	data, err := json.Marshal(p.Database.GetBookFromDatabase(idint))
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(data)
}

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

	result := p.Database.SaveBookToDataBase(newBook)

	data, err := json.Marshal(result)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(data)
}

func (p Server) AllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	limit := query.Get("limit")

	if limit == "" {

		data, err := json.Marshal(p.Database.GetAllBookFromDatabase())
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		w.Write(data)
		Logger.Info("отправлен ответ")
		return
	}
	if limit != "" {
		limitNum, err := strconv.Atoi(limit)
		if err != nil {
			handleError(w, http.StatusBadRequest, errors.New("invalid limit parameter"))
			return
		}

		if limitNum > len(p.Database.Store) {
			limitNum = len(p.Database.Store)
		}

		for i := 1; i <= limitNum; i++ {
			data, err := json.Marshal(p.Database.GetBookFromDatabase(i))
			if err != nil {
				handleError(w, http.StatusInternalServerError, err)
				return
			}
			w.Write(data)
			Logger.Info("отправлен ответ")
		}
		return
	}
}

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
	p.Database.UpDateBookToDataBase(book, book.Id)
}

func (p Server) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idstr := r.URL.Query().Get("id")
	idint, err := strconv.Atoi(idstr)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	p.Database.DeleteBookFromDatabase(idint)
	w.WriteHeader(http.StatusNoContent)
	Logger.Info("удалена книга")
}
