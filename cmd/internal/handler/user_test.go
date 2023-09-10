package handler

import (
	"ebook/cmd/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.UserService

		reqBody  string
		wantCode int
		wantBody string
	}{
		{},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}

func TestNil(t *testing.T) {
	testTypeAssert(nil)
}

func testTypeAssert(c any) {
	_, ok := c.(*UserClaims)
	println(ok)
}

func TestEncrypt(t *testing.T) {
	_ = NewUserHandler(nil, nil)
	password := "hello#world123"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)
}
