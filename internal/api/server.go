package api

import (
	pb "books1/internal/api/proto/v1"
	"books1/internal/domain"
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
	//Database *db.Repository `json:"database"`
	Database Store
}

func (p Server) GetBook(ctx context.Context, request *pb.GetBookRequest) (*pb.GetBookResponse, error) {

	idint := uint(request.Id)
	result, err := p.Database.GetBookFromDatabaseByRAWSql(ctx, uint(idint))
	if err != nil {
		fmt.Println("error10")
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

/*
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
*/
func (p Server) AllBooks(ctx context.Context, r *pb.AllBooksRequest) (*pb.AllBooksResponse, error) {

	books, err := p.Database.GetAllBookFromDatabaseByRAWSql(ctx)

	if err != nil {
		fmt.Println(err)
	}

	g := []*pb.Book1{}
	n := &pb.Book1{}

	for i := 0; i < len(books); i++ {

		n = &pb.Book1{
			Id:    int64(books[i].ID),
			Title: books[i].Title,
			Year:  int32(books[i].Year),
		}
		g = append(g, n)
	}
	return &pb.AllBooksResponse{Book1: g}, nil
}

/*
books := []domain.Book{}

	if err != nil {
		fmt.Println("error3")
	}
	for rows.Next() {
		book := &domain.Book{}
		rows.Scan(&book.ID, &book.Title, &book.Year)
		books = append(books, *book)
	}
*/
//return &pb.AllBooksResponse{Book: &pb.Book{
//	Id:    int64(books.ID),
//	Title: books.Title,
//	Year:  int32(books.Year),
//}}, nil

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
