package api

import (
	pb "books1/internal/api/proto/v1"
	"books1/internal/domain"
	"books1/internal/logger"
	"context"
	"fmt"
)

type Store interface {
	SaveBookToDataBaseByRAWSql(ctx context.Context, book domain.Book) (*domain.Book, error)
	GetBookFromDatabaseByRAWSql(ctx context.Context, id uint) (*domain.Book, error)
	GetAllBookFromDatabaseByRAWSql(ctx context.Context) ([]domain.Book, error)
	DeleteBookFromDatabaseByRAWSql(ctx context.Context, id uint) error
	UpDateBookToDataBaseByRAWSql(ctx context.Context, book domain.Book) error
}

type Server struct {
	Database Store
}

func (p Server) GetBook(ctx context.Context, request *pb.GetBookRequest) (*pb.GetBookResponse, error) {

	idint := uint(request.Id)
	result, err := p.Database.GetBookFromDatabaseByRAWSql(ctx, uint(idint))
	if err != nil {
		fmt.Println("error10")
	}

	log, found := logger.FromContext(ctx)
	if found == false {
		log.Debug("Проблемы у нас")

	}
	log.Info("отправляем в ответ книгу")

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
	log, found := logger.FromContext(ctx)
	if found == false {
		log.Debug("Проблемы у нас")

	}
	log.Info("добавили книгу")
	return &pb.AddBookResponse{Book: &pb.Book{
		Id:    int64(result.ID),
		Title: result.Title,
		Year:  int32(result.Year),
	}}, nil

}

func (p Server) AllBooks(ctx context.Context, r *pb.AllBooksRequest) (*pb.AllBooksResponse, error) {

	books, err := p.Database.GetAllBookFromDatabaseByRAWSql(ctx)

	if err != nil {
		fmt.Println(err)
	}

	a := int(len(books))
	if r.Limit > 0 {
		a = int(r.Limit)
	}

	g := []*pb.Book1{}
	n := &pb.Book1{}

	for i := 0; i < a; i++ {

		n = &pb.Book1{
			Id:    int64(books[i].ID),
			Title: books[i].Title,
			Year:  int32(books[i].Year),
		}
		g = append(g, n)
	}
	log, found := logger.FromContext(ctx)
	if found == false {
		log.Debug("Проблемы у нас")

	}
	log.Info("отправляем в ответ книги")
	return &pb.AllBooksResponse{Book1: g}, nil
}

func (p Server) UpdateBook(ctx context.Context, request *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {

	newBook := domain.Book{
		Title: request.Title,
		Year:  int(request.Year),
	}

	err := p.Database.UpDateBookToDataBaseByRAWSql(ctx, newBook)
	if err != nil {
		fmt.Println("книга", request.Id, "обновлена")
	}
	log, found := logger.FromContext(ctx)
	if found == false {
		log.Debug("Проблемы у нас")

	}
	log.Info("обновили данные книги")
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
		fmt.Println("проблема, книга не удалена")

	}
	log, found := logger.FromContext(ctx)
	if found == false {
		log.Debug("Проблемы у нас")

	}
	log.Info("книга удалена")

	return &pb.DeleteBookResponse{
		Id: int64(request.Id),
	}, nil

}
