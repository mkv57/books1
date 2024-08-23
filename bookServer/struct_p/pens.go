package struct_p

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idstr := r.URL.Query().Get("id")
	idint, err := strconv.Atoi(idstr)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	book, ok := Books[idint]
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

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsong, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	var newBook Book
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
	newBook.Id = len(Books) + 1 // формируем новый идентификатор
	Books[len(Books)+1] = newBook
	data, err := json.Marshal(newBook)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(data)
	logger.Info("добавлена книга")
}

func AllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	limit := query.Get("limit")

	if limit == "" {

		data, err := json.Marshal(Books)
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
		if limitNum > len(Books) {
			limitNum = len(Books)
		}

		for i := 1; i <= limitNum; i++ {
			data, err := json.Marshal(Books[i])
			if err != nil {
				handleError(w, http.StatusInternalServerError, err)

			}
			w.Write(data)
			logger.Info("отправлен ответ")
		}
		return
	}

}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	var book Book
	err = json.Unmarshal(data, &book)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	Books[book.Id] = book
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idstr := r.URL.Query().Get("id")
	idint, err := strconv.Atoi(idstr)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}
	delete(Books, idint)
	w.WriteHeader(http.StatusNoContent)
	logger.Info("удалена книга")
}

func handleError(w http.ResponseWriter, status int, err error) {
	result := map[string]interface{}{
		"error":  err.Error(),
		"status": http.StatusText(status),
	}
	data, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
}
