package api

import (
	"books1/internal/domain"
	"context"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestServer_GetBook(t *testing.T) {
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
