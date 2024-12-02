package api

import (
	"books1/internal/domain"
	"books1/internal/logger"
	"context"
	"encoding/json"
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

	expectedBook := &domain.Book{
		ID:        1,
		Title:     "ttttttt",
		Year:      2021,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx := context.Background()

	ctx = logger.NewContext(ctx, log)

	mockStore.EXPECT().
		GetBookFromDatabaseByRAWSql(ctx, uint(1)).
		Return(expectedBook, nil)

	reg, err := http.NewRequest(http.MethodGet, "http://localhost:8080/book?id=1", nil)

	require.NoError(t, err) // проверка используется вместо if t!= nil  но только для тестов

	reg = reg.WithContext(ctx)

	resp := httptest.NewRecorder()

	server.GetBook(resp, reg)

	result := &domain.Book{}

	err = json.Unmarshal(resp.Body.Bytes(), &result)
	require.NoError(t, err)

	expectedBook.CreatedAt = result.CreatedAt
	expectedBook.UpdatedAt = result.UpdatedAt

	require.Equal(t, expectedBook, result)
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
