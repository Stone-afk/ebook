package integration

import (
	"database/sql"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/handler"
	"ebook/cmd/internal/integration/startup"
	dao "ebook/cmd/internal/repository/dao/user"
	"ebook/cmd/internal/service/oauth2"
	oauth2mocks "ebook/cmd/internal/service/oauth2/mocks"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 这是一个只能手动运行的测试，为了摆脱 wechat 那个部分而引入的测试
func TestWechatCallback(t *testing.T) {
	const callbackUrl = "/oauth2/wechat/callback"
	db := startup.InitTestDB()
	testCases := []struct {
		name   string
		mock   func(ctrl *gomock.Controller) oauth2.Service
		before func(t *testing.T)
		// 验证并且删除数据
		after      func(t *testing.T)
		wantCode   int
		wantResult handler.Result
	}{
		{
			name: "注册新用户",
			mock: func(ctrl *gomock.Controller) oauth2.Service {
				svc := oauth2mocks.NewMockService(ctrl)
				svc.EXPECT().VerifyCode(gomock.Any(),
					gomock.Any()).
					Return(domain.WechatInfo{
						OpenID:  "123",
						UnionID: "1234",
					}, nil)
				return svc
			},
			before: func(t *testing.T) {
				// 什么也不需要做
			},
			after: func(t *testing.T) {
				// 验证数据库
				var u dao.User
				err := db.Find(&u, "wechat_open_id = ?", "123").Error
				assert.NoError(t, err)
				// 只需要验证 union id 就差不多了
				assert.Equal(t, "1234", u.WechatUnionID.String)
				db.Delete(&u, "wechat_open_id = ?", "123")
			},
			wantCode: 200,
			wantResult: handler.Result{
				Msg: "登录成功",
			},
		},
		{
			name: "已有的用户",
			mock: func(ctrl *gomock.Controller) oauth2.Service {
				svc := oauth2mocks.NewMockService(ctrl)
				svc.EXPECT().VerifyCode(gomock.Any(), gomock.Any()).
					Return(domain.WechatInfo{
						OpenID:  "2345",
						UnionID: "23456",
					}, nil)
				return svc
			},
			before: func(t *testing.T) {
				// 插入数据，假装用户存在
				err := db.Create(&dao.User{
					WechatOpenID: sql.NullString{
						String: "2345",
						Valid:  true,
					},
					WechatUnionID: sql.NullString{
						String: "23456",
						Valid:  true,
					},
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据库
				var u dao.User
				err := db.Find(&u, "wechat_open_id = ?", "2345").Error
				assert.NoError(t, err)
				// 只需要验证 union id 就差不多了
				assert.Equal(t, "23456", u.WechatUnionID.String)
				db.Delete(&u, "wechat_open_id = ?", "123")
			},
			wantCode: 200,
			wantResult: handler.Result{
				Msg: "登录成功",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tc.before(t)

			userSvc := startup.InitUserSvc()
			jwtHdl := startup.InitJwtHandler()
			wechatSvc := tc.mock(ctrl)
			hdl := handler.NewOAuth2WechatHandler(wechatSvc, userSvc,
				jwtHdl, startup.NewWechatHandlerConfig())
			server := gin.Default()
			hdl.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodGet,
				callbackUrl, nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			code := recorder.Code
			// 反序列化为结果
			var result handler.Result
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, code)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}
