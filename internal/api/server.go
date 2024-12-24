package api

import (
	pb "books1/internal/api/proto/v1"
	"books1/internal/domain"
	"books1/internal/logger"
	"context"
	"encoding/json"
	"fmt"

	"errors"
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

func (p Server) GetBook(ctx context.Context, request *pb.GetBookRequest) (*pb.GetBookResponse, error) {

	idint := uint(request.Id)
	result, err := p.Database.GetBookFromDatabaseByRAWSql(ctx, uint(idint))
	if err != nil {
		fmt.Println("error10") //handleError(w, http.StatusInternalServerError, err)
	}
	return &pb.GetBookResponse{Book: &pb.Book{
		Id:    int64(result.ID),
		Title: result.Title,
		Year:  int32(result.Year),
	}}, nil

}

func (p Server) AddBook(ctx context.Context, request *pb.AddBookRequest) (*pb.AddBookResponse, error) {

	newBook := domain.Book{
		Title: request.Title,
		Year:  int(request.Year),
	}

	result, err := p.Database.SaveBookToDataBaseByRAWSql(ctx, newBook)
	if err != nil {
		return nil, err
	}
	fmt.Println("книга добавлена")
	return &pb.AddBookResponse{Book: &pb.Book{
		Id:    int64(result.ID),
		Title: result.Title,
		Year:  int32(result.Year),
	}}, nil

}

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

/*
func (p Server) AllBooks(ctx context.Context, r *pb.AllBooksRequest) (*pb.AllBooksResponse, error) {

		//ctx := r.Context()
		//query := r.URL.Query()
		//limit := query.Get("limit")

		books, err := p.Database.GetAllBookFromDatabaseByRAWSql(ctx)
		//fmt.Println(books)
		if err != nil {
			//	handleError(w, http.StatusInternalServerError, err)
			fmt.Println(err)
		}
		return books, nil
		//if limit != "" {
		//	limitNum, err := strconv.Atoi(limit)
		//	if err != nil {
		//		handleError(w, http.StatusBadRequest, errors.New("invalid limit parameter"))
		//		return
		//	}
		//	fmt.Println(limitNum)
		//}

		//data, err := json.Marshal(books)
		//if err != nil {
		//	handleError(w, http.StatusInternalServerError, err)
		//	return
		//}

		//w.Write(data)

		//log, found := logger.FromContext(ctx)
		//if found == false {
		//	handleError(w, http.StatusInternalServerError, errors.New("Проблемы у нас"))
		//	return
		//}
		//log.Info("отправлен ответ")
	}
*/

func (p Server) UpdateBook(ctx context.Context, request *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {

	newBook := domain.Book{
		Title: request.Title,
		Year:  int(request.Year),
	}

	err := p.Database.UpDateBookToDataBaseByRAWSql(ctx, newBook)
	if err != nil {
		fmt.Println("книга", request.Id, "обновлена")
	}
	return &pb.UpdateBookResponse{Book: &pb.Book{
		Id:    int64(request.Id),
		Title: request.Title,
		Year:  int32(request.Year),
	}}, nil
}

func (p Server) DeleteBook(ctx context.Context, request *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	idint := uint(request.Id)

	err := p.Database.DeleteBookFromDatabaseByRAWSql(ctx, idint)
	if err != nil {
		fmt.Println("проблема, книга не удалена") // handleError(w, http.StatusBadRequest, err)

	}
	fmt.Println("книга удалена")

	return &pb.DeleteBookResponse{
		Id: int64(request.Id),
	}, nil

}
