package api

import (
	pb "books1/internal/api/proto/v1"
	"books1/internal/domain"
	"books1/internal/logger"
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Store interface {
	SaveBookToDataBaseByRAWSql(ctx context.Context, book domain.Book) (*domain.Book, error)
	GetBookFromDatabaseByRAWSql(ctx context.Context, id uint) (*domain.Book, error)
	GetAllBookFromDatabaseByRAWSql(ctx context.Context) ([]domain.Book, error)
	DeleteBookFromDatabaseByRAWSql(ctx context.Context, id uint) error
	UpDateBookToDataBaseByRAWSql(ctx context.Context, book domain.Book) error
	// User logic
	SaveUserToDatabase(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	SaveSessionTodatabase(ctx context.Context, session domain.Session) error
	GetUserByToken(ctx context.Context, token string) (domain.User, error)
}

type Server struct {
	Database Store
}

// Login implements pb.BookAPIServer.
var (
	ErrInvalidPassword = errors.New("invalid password")
)

func (p *Server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := p.Database.GetUserByEmail(ctx, request.Email) // ищем данные user password  и id по email
	if err != nil {
		return nil, err
	}
	if user.Password != request.Password {
		return nil, ErrInvalidPassword
	}

	ip, err := originFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	token := uuid.New().String()
	session := domain.Session{
		UserID:    user.ID,
		Token:     token,
		IP:        ip,
		UserAgent: "", // TODO

	}
	err = p.Database.SaveSessionTodatabase(ctx, session)
	if err != nil {
		return nil, err
	}
	err = grpc.SendHeader(ctx, metadata.MD{"authorization": {token}})
	if err != nil {
		return nil, fmt.Errorf("error")
	}
	return &pb.LoginResponse{User: &pb.User{
		Id:    int64(user.ID),
		Email: user.Email,
	}}, nil
}

// Функция достаёт IP клиента
func originFromCtx(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("error")
	}
	p.Addr.String()

	clientIP, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return "", fmt.Errorf("error")
	}
	return clientIP, nil
}

// Registration implements pb.BookAPIServer.
func (p *Server) Registration(ctx context.Context, request *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	newUser := domain.User{
		Email:    request.Email,
		Password: request.Password,
	}
	registeredUser, err := p.Database.SaveUserToDatabase(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return &pb.RegistrationResponse{Id: int64(registeredUser.ID)}, nil
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

const authScheme = "Bearer"

func (p Server) AddBook(ctx context.Context, request *pb.AddBookRequest) (*pb.AddBookResponse, error) {

	token, err := auth.AuthFromMD(ctx, authScheme)
	if err != nil {
		return nil, fmt.Errorf("")
	}

	user, err := p.Database.GetUserByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	newBook := domain.Book{
		Title:  request.Title,
		Year:   int(request.Year),
		UserID: user.ID,
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
