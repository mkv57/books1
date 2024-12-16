package api

import (
	"books1/internal/domain"
	"books1/internal/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"log/slog"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	//"golang.org/x/exp/slog"
)

func TestServer_GetBook(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockStore := NewMockStore(ctrl)

	server := Server{

		Database: mockStore,
	}

	expectedBook := &domain.Book{ // создаём экземпляр структуры Book
		ID:        1,
		Title:     "ttttttt",
		Year:      2021,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ // создаём переменную log для добавления в context
		Level: slog.LevelDebug,
	}))

	ctx := context.Background() // создаём переменную типа context с пустым значением Background

	ctx = logger.NewContext(ctx, log) // вкладываем в context 2е переменные ctx и log , как и в функции GetBookFromDatabaseByRAWSql

	mockStore.EXPECT().
		GetBookFromDatabaseByRAWSql(ctx, uint(1)).
		Return(expectedBook, nil) // ???

	reg, err := http.NewRequest(http.MethodGet, "http://localhost:8080/book?id=1", nil) // оправляем запрос на обработчик GetBook, который вызывает
	// функцию GetBookFromDatabaseByRAWSql(ctx, uint(idint)), он достаёт  данные из структуры Book и эти данные отправляются
	// в Request reg	???

	require.NoError(t, err) // проверка используется вместо if t!= nil  но только для тестов

	reg = reg.WithContext(ctx) // добавляем в запрос в context то же значение ctx

	resp := httptest.NewRecorder() // ??? создаём переменную типа ResponseWrite, чтобы вызвать GetBook, вложив resp
	fmt.Println("1", resp.Body, "q", reg)

	server.GetBook(resp, reg) // вызываем функцию, которая должна записать в тело ответа resp.Body  json c данными структуры Book
	fmt.Println("2", resp.Body, "q", reg)
	result := &domain.Book{} // создаём переменную типа структуры Book

	err = json.Unmarshal(resp.Body.Bytes(), &result) // resp.Body.Bytes() что это ???	распарсили json и вложили данные в result
	require.NoError(t, err)                          // проверка на ошибки

	expectedBook.CreatedAt = result.CreatedAt // уровняли зачения времени, так как будут разные значения
	expectedBook.UpdatedAt = result.UpdatedAt // уровняли зачения времени

	require.Equal(t, expectedBook, result) // сравнили результаты
}

/*func TestServer_GetBook(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockStore := NewMockStore(ctrl)
	srv := Server{

		Database: mockStore,
	}

	mockStore.EXPECT().GetBookFromDatabaseByRAWSql(context.Background(), uint(1)).
		Return(&domain.Book{
			ID:        1,
			Title:     "ttttttt",
			Year:      2021,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}, nil)

	mockStore.EXPECT().GetAllBookFromDatabaseByRAWSql(context.Background()).
		Return([]domain.Book{
			{
				ID:        1,
				Title:     "ttttttt",
				Year:      2021,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		})

	srv.GetBook(nil, nil)
}
*/
