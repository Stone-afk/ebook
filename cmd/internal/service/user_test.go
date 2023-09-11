package service

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("hello#world123"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}

func Test_userService_Login(t *testing.T) {
	// 做成一个测试用例都用到的时间
	// now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 输入
		//ctx      context.Context
		email    string
		password string

		// 输出
		wantUser domain.User
		wantErr  error
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// 具体的测试代码
			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)

		})
	}
}
