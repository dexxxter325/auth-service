package handler

import (
	"CRUD_API"
	"CRUD_API/pkg/service"
	mockservice "CRUD_API/pkg/service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestALl(t *testing.T) {
	t.Run("Create", TestHandler_CreateProduct)
	t.Run("ReadAll", TestHandler_ReadAllProducts)
	t.Run("ReadById", TestHandler_ReadProductById)
	t.Run("Update", TestHandler_UpdateProduct)
	t.Run("Delete", TestHandler_DeleteProduct)
}
func clearCache() {
	// Очистка кэша
	cache := GetSingleton()
	cache.Clear()
}

func TestHandler_CreateProduct(t *testing.T) {
	type Mock func(s *mockservice.MockProduct, name, description string) //наши моки<>c *mock_handler.MockCaches
	/*createdProduct := CRUD_API.Products{
		ID:          1,
		Name:        "ok",
		Description: "ok",
	}*/
	type args struct {
		name        string
		description string
	}
	tests := []struct {
		name               string
		request            string //наш запрос в postman
		mock               Mock
		args               args
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "ok",
			args: args{
				name:        "ok",
				description: "ok",
			},
			request: `{"name":"ok","description":"ok"}`,
			mock: func(s *mockservice.MockProduct, name, description string) {
				s.EXPECT().Create(name, description).Return(CRUD_API.Products{ //expect-моковый вызов функции create
					ID:          1,
					Name:        name,
					Description: description,
				}, nil)
				//c.EXPECT().Set("product:1", createdProduct, gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode: 200,
			expectedResponse:   `{"product(created successfully)":{"id":1,"name":"ok","description":"ok"}}`,
		},
		{
			name:               "invalid json",
			request:            `{"invalid json man!"}`,
			mock:               func(s *mockservice.MockProduct, name, description string) {},
			expectedStatusCode: 400,
			expectedResponse:   `{"error":"Invalid JSON in CreateProduct"}`,
		},
		{
			name:    "service err",
			request: `{"name":"test","description":"test"}`,
			args: args{
				name:        "test",
				description: "test",
			},
			mock: func(s *mockservice.MockProduct, name, description string) {
				s.EXPECT().Create(name, description).Return(CRUD_API.Products{}, errors.New("err in service"))
			},
			expectedStatusCode: 500,
			expectedResponse:   `{"error":"err in service"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			product := mockservice.NewMockProduct(ctrl)
			//cache := mockcache.NewMockCaches(ctrl)
			tt.mock(product, tt.args.name, tt.args.description)
			services := &service.Service{Product: product}
			handler := Handler{services}
			r := gin.New()
			r.POST("/product", handler.CreateProduct)
			req := httptest.NewRequest("POST", "/product", bytes.NewBufferString(tt.request))
			/*bytes.NewBufferString преобразовывает наше тело запроса к io.reader*/
			w := httptest.NewRecorder() //ответ
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
			clearCache()
		})
	}
}
func TestHandler_ReadAllProducts(t *testing.T) {
	type Mock func(s *mockservice.MockProduct)
	tests := []struct {
		name               string
		mock               Mock
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "ok",
			mock: func(s *mockservice.MockProduct) {
				expectedProducts := []CRUD_API.Products{
					{ID: 1, Name: "ok", Description: "ok"},
				}
				s.EXPECT().ReadAll().Return(expectedProducts, nil)
			},
			expectedStatusCode: 200,
			expectedResponse:   `{"product":[{"id":1,"name":"ok","description":"ok"}]}`,
		},
		{
			name: "err in service",
			mock: func(s *mockservice.MockProduct) {
				s.EXPECT().ReadAll().Return([]CRUD_API.Products{}, errors.New("err in service"))
			},
			expectedStatusCode: 500,
			expectedResponse:   `{"error":"err in service"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			product := mockservice.NewMockProduct(ctrl)
			tt.mock(product)
			services := &service.Service{Product: product}
			handler := Handler{services}
			r := gin.New()
			r.GET("/product", handler.ReadAllProducts)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/product", nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
			assert.Equal(t, tt.expectedStatusCode, w.Code)

		})
	}
}
func TestHandler_ReadProductById(t *testing.T) {
	//GetSingleton()
	type Mock func(s *mockservice.MockProduct, id int)
	tests := []struct {
		name               string
		id                 int
		mock               Mock
		expectedStatusCode int
		expectedResponse   string
	}{
		/*{
			name: "ok from cache",
			id:   1,
			mock: func(s *mockservice.MockProduct, c *mockcache.MockCaches, id int) {
				expectedProduct := CRUD_API.Products{ID: id, Name: "ok", Description: "ok"}
				c.EXPECT().Set("product:1", gomock.Any(), cacheTTl)
				c.EXPECT().Get("product:1").Do(func(key string) {
					fmt.Printf("cacheMock.Get(%s) called\n", key)
				}).Return(expectedProduct)
				s.EXPECT().ReadById(id).Return(expectedProduct, nil)

				//cache.On("Get", "product:1").Return(expectedProduct)
			},
			expectedResponse:   `{"from_cache":true,"product":{"id":1,"name":"ok","description":"ok"}}`,
			expectedStatusCode: 200,
		},*/
		{
			name: "ok from db",
			id:   1,
			mock: func(s *mockservice.MockProduct, id int) {
				expectedProduct := CRUD_API.Products{ID: id, Name: "ok", Description: "ok"}
				s.EXPECT().ReadById(id).Return(expectedProduct, nil)
				//cache.On("Get", "product:1").Return(nil) // Мок для чтения из кэша
			},
			expectedResponse:   `{"from_cache":false,"product":{"id":1,"name":"ok","description":"ok"}}`,
			expectedStatusCode: 200,
		},
		{
			name: "err in service",
			id:   777,
			mock: func(s *mockservice.MockProduct, id int) {
				s.EXPECT().ReadById(id).Return(CRUD_API.Products{}, errors.New("err in service"))
			},
			expectedResponse:   `{"error":"err in service"}`,
			expectedStatusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			product := mockservice.NewMockProduct(ctrl)
			//cache := mockcache.NewMockCaches(ctrl)
			tt.mock(product, tt.id)
			services := &service.Service{Product: product}
			handler := Handler{services}
			r := gin.New()
			r.GET("/product/:id", handler.ReadProductById)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/product/%d", tt.id), nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_UpdateProduct(t *testing.T) {
	mockedRedisClient, mock := redismock.NewClientMock()
	RedisClient = mockedRedisClient
	product := CRUD_API.Products{
		ID:          1,
		Name:        "ok",
		Description: "ok",
	}
	MarshalForRedis, _ := json.Marshal(product)
	type Mock func(s *mockservice.MockProduct, name, description string, id int)
	type args struct {
		name        string
		description string
	}
	tests := []struct {
		id                 int
		name               string
		request            string
		args               args
		mock               Mock
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "ok",
			id:   1,
			args: args{
				name:        "ok",
				description: "ok",
			},
			request: `{"id":1,"name":"ok","description":"ok"}`,
			mock: func(s *mockservice.MockProduct, name, description string, id int) {
				updatedProduct := CRUD_API.Products{Name: name, Description: description, ID: id}
				s.EXPECT().Update(name, description, id).Return(updatedProduct, nil)
				mock.ExpectSet("product:1", MarshalForRedis, cacheTTl).SetVal("nil")
			},
			expectedResponse:   `{"product(updated successfully)":{"id":1,"name":"ok","description":"ok"}}`,
			expectedStatusCode: 200,
		},
		{
			name:               "invalid json",
			request:            `"invalid json"`,
			mock:               func(s *mockservice.MockProduct, name, description string, id int) {},
			expectedResponse:   `{"error":"invalid json in UpdateProduct"}`,
			expectedStatusCode: 400,
		},
		{
			name: "service err",
			id:   1,
			args: args{
				name:        "test",
				description: "test",
			},
			request: `{"id":1,"name":"test","description":"test"}`,
			mock: func(s *mockservice.MockProduct, name, description string, id int) {
				updatedProduct := CRUD_API.Products{Name: name, Description: description, ID: id}
				s.EXPECT().Update(name, description, id).Return(updatedProduct, errors.New("err in service"))
			},
			expectedStatusCode: 500,
			expectedResponse:   `{"error marshaling product to JSON in update":"err in service"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			product := mockservice.NewMockProduct(ctrl)
			tt.mock(product, tt.args.name, tt.args.description, tt.id)
			services := &service.Service{Product: product}
			handler := Handler{services}
			r := gin.New()
			r.PUT("/product/:id", handler.UpdateProduct)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/product/%d", tt.id), bytes.NewBufferString(tt.request))
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
func TestHandler_DeleteProduct(t *testing.T) {
	mockedRedisClient, mock := redismock.NewClientMock()
	RedisClient = mockedRedisClient
	type Mock func(s *mockservice.MockProduct, id int)
	tests := []struct {
		name               string
		mock               Mock
		id                 int
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "ok",
			id:   1,
			mock: func(s *mockservice.MockProduct, id int) {
				s.EXPECT().Delete(id).Return(nil)
				mock.ExpectDel("product:1").SetVal(0)
			},
			expectedStatusCode: 200,
			expectedResponse:   `"deleted successfully"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			product := mockservice.NewMockProduct(ctrl)
			tt.mock(product, tt.id)
			services := &service.Service{Product: product}
			handler := Handler{services}
			r := gin.New()
			r.DELETE("/product/:id", handler.DeleteProduct)
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/product/%d", tt.id), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
		})
	}
}
