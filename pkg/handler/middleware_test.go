package handler

import (
	"CRUD_API/pkg/service"
	mockservice "CRUD_API/pkg/service/mocks"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_userIdentity(t *testing.T) {
	type Mock func(s *mockservice.MockAuthorization, token string)

	tests := []struct {
		name               string
		token              string
		headerName         string
		headerValue        string
		mock               Mock
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mock: func(s *mockservice.MockAuthorization, token string) {
				s.EXPECT().ParseAccessToken(token).Return(1, nil)
			},
			expectedResponse:   "1",
			expectedStatusCode: 200,
		},
		{
			name:               "no header",
			headerName:         "",
			headerValue:        "Bearer token",
			token:              "token",
			mock:               func(s *mockservice.MockAuthorization, token string) {},
			expectedResponse:   `{"error":"empty auth header"}`,
			expectedStatusCode: 401,
		},
		{
			name:               "wrong auth header",
			headerName:         "Authorization",
			headerValue:        "BeaRRer token",
			token:              "token",
			mock:               func(s *mockservice.MockAuthorization, token string) {},
			expectedResponse:   `{"error":"invalid auth header"}`,
			expectedStatusCode: 401,
		},
		{
			name:               "no token",
			headerName:         "Authorization",
			headerValue:        "Bearer ",
			mock:               func(s *mockservice.MockAuthorization, token string) {},
			expectedResponse:   `{"error":"empty token"}`,
			expectedStatusCode: 401,
		},
		{
			name:        "err in parse",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mock: func(s *mockservice.MockAuthorization, token string) {
				s.EXPECT().ParseAccessToken(token).Return(0, errors.New("failed in parse token"))
			},
			expectedResponse:   `{"error in parse":"failed in parse token"}`,
			expectedStatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			auth := mockservice.NewMockAuthorization(ctrl)
			tt.mock(auth, tt.token)
			services := &service.Service{Authorization: auth}
			handler := Handler{services}
			r := gin.New()
			r.GET("/product", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get("userId")
				c.String(200, "%d", id) //если id получили,то выводим успешное завершение
			})
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/product", nil)
			req.Header.Set(tt.headerName, tt.headerValue) //передаем наши значения в заголовке в запросе
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Body.String(), tt.expectedResponse)
			assert.Equal(t, w.Code, tt.expectedStatusCode)

		})
	}
}
