package service

import (
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_articleService_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository,
			repository.ArticleReaderRepository)

		art domain.Article

		wantErr error
		wantId  int64
	}{
		{},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}
