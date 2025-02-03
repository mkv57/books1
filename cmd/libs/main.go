package main

//dsт: "host=localhost user=mkv password=book_server dbname=book_database port=5432 sslmode=disable"
//dsn2: "postgres://mkv:book_server@localhost:5435/book_database?sslmode=disable"

import (
	"context"
	"database/sql"
	"log"
	"net"

	grpc_gateway_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"

	"books1/internal/api"
	pb "books1/internal/api/proto/v1"
	"books1/internal/db"
	"books1/internal/logger"

	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"

	_ "github.com/lib/pq"
)

type Config struct {
	DSN          string `yaml:"dsn"`
	LogLevel     int    `yaml:"log_level"`
	Address_grps string `yaml:"address_grps"`
	Address_http string `yaml:"address_http"`
}

func main() {

	yamlContent, err := os.ReadFile("./config.yml") // читаем файл config.yml и сохраняем в переменную yamlContent
	if err != nil {
		log.Fatal(err)
	}
	var systemconfig Config
	err = yaml.Unmarshal(yamlContent, &systemconfig) // декодировали данные из переменной yamlContent в переменную systemconfig
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile("./app.log", os.O_APPEND, 0666) // Открываем файл app.log для записи туда логов
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var loggerOur = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ // создаём переменную, для передачи логов в формате
		Level: slog.Level(systemconfig.LogLevel), // (loggerOur.Warn("сервер запущен")) Если os.Stdout, то в терминал
	})) // Если переменная file то в файл app.log

	m, err := migrate.New( // сохраняем в переменную m данные из папки migrate
		"file://migrate", systemconfig.DSN)
	if err != nil {

		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange { // отправляем запрос из переменной m в базу данных sql из папки migrate
		fmt.Println("миграция не прошла")
		log.Fatal(err)
	}

	rawSQLConn, err := sql.Open("postgres", systemconfig.DSN) // открываем доступ к базе данных postgreSQL
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//r := mux.NewRouter()

	repo := db.NewRepository(rawSQLConn) // зоздаём путь для передачи и получения данных из базы данных postgreSQL

	//r.Use(api.Log(loggerOur))

	ourServer := api.Server{ // создаём связь(путь) передачи данных между обработчиками и базой данной postgreSQL

		Database: repo,
	}

	// создаём структуру которая слушает порты для получениязапроса по grpc( и только)
	ln, err := net.Listen("tcp", systemconfig.Address_grps)
	if err != nil {
		log.Fatal("net.Listen(): %v", err)
	}

	server := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(validator.UnaryServerInterceptor()),
		grpc.ChainUnaryInterceptor(
			loginterceptor(loggerOur),
			logging.UnaryServerInterceptor(interceptorLogger(loggerOur)),
		),
	)
	// здесь регестрируем хендеры для grpc
	pb.RegisterBookAPIServer(server, &ourServer)
	// здесь запускаем только grpc сервер
	go func() {
		if err := server.Serve(ln); err != nil {
			log.Fatal("server.Serve(): %v", err)
		}
	}()
	// здесь создаём конструкцию, которая умеет вызывать (grpc) методы
	conn, err := grpc.NewClient(systemconfig.Address_grps,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("grpc.NewClient(): %v", err)
	}
	defer conn.Close()
	// TODO: fix me
	// HTTP  сервер не умеет обрабатывать grpc запросы, поэтому используем grpc-gateway
	var gw = grpc_gateway_runtime.NewServeMux() // нужно сделать тег grpc_gateway_runtime сдесь в main.go
	// grpc- gateway принимает клиент до grpc сервера и преобразует HTTP (REST) в grpc
	// и перенаправляет через client
	err = pb.RegisterBookAPIHandler(context.Background(), gw, conn)
	if err != nil {
		log.Fatal("RegisterServiceAPIHandler(): %v", err)
	}
	// web стучится по REST в gateway, который конвертирует из HTTP в grpc
	// тот получил и отправил grpc-client и тот в grpc-server
	gwServer := &http.Server{
		Addr:    systemconfig.Address_http,
		Handler: gw,
	}

	loggerOur.Warn("сервер запущен")
	err1 := gwServer.ListenAndServe()
	//err = http.ListenAndServe("localhost:8080", r)
	loggerOur.Warn("сервер отключён")
	if err1 != nil {
		fmt.Println(err1)
		loggerOur.Debug("сервер нe запустился")
	}

}
func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func loginterceptor(loggerOur *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, reg any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, er error) {
		log1 := loggerOur.With(
			slog.String("full_method", info.FullMethod),
		)
		ctx = logger.NewContext(ctx, log1)
		return handler(ctx, reg)
	}
}
