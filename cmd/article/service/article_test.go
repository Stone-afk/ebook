package service

import (
	"context"
	"ebook/cmd/article/domain"
	"ebook/cmd/article/repository"
	repomocks "ebook/cmd/article/repository/mocks"
	"ebook/cmd/pkg/logger/nop"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_articleService_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.ArticleRepository

		art domain.Article

		wantErr error
		wantId  int64
	}{
		{
			name: "新建发表成功",
			mock: func(ctrl *gomock.Controller) repository.ArticleRepository {
				repo := repomocks.NewMockArticleRepository(ctrl)
				repo.EXPECT().Sync(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return repo
			},
			art: domain.Article{
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 1,
		},
		{
			name: "修改并发表成功",
			mock: func(ctrl *gomock.Controller) repository.ArticleRepository {
				repo := repomocks.NewMockArticleRepository(ctrl)
				repo.EXPECT().Sync(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				return repo
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 2,
		},
		{
			// 新建-保存到制作库失败
			// 修改-保存到制作库失败
			name: "保存到制作库失败",
			mock: func(ctrl *gomock.Controller) repository.ArticleRepository {
				repo := repomocks.NewMockArticleRepository(ctrl)
				repo.EXPECT().Sync(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("mock db error"))
				return repo
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  0,
			wantErr: errors.New("mock db error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			svc := NewArticleService(repo, nil, nil, nil, nop.NewNopLogger())
			id, err := svc.Publish(context.Background(), tc.art)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, id)
		})
	}
}

func Test_articleService_PublishV1(t *testing.T) {
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			author, reader := tc.mock(ctrl)
			svc := NewArticleServiceV1(author, reader, nil, nop.NewNopLogger())
			id, err := svc.PublishV1(context.Background(), tc.art)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, id)
		})
	}
}
