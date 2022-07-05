package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/service"
	mock_service "github.com/AndrewMislyuk/go-shop-backend/internal/service/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_createProduct(t *testing.T) {
	type mockBehavior func(s *mock_service.MockProductsList, input domain.CreateProductInput)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           domain.CreateProductInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"title":"test_title","image":"test_image","price":70000,"sale":0,"sale_old_price":0,"category":"test_category","type":"test_type","subtype":"test_subtype","description":"test_description"}`,
			inputUser: domain.CreateProductInput{
				Title:        "test_title",
				Price:        70000,
				Sale:         0,
				SaleOldPrice: 0,
				Category:     "test_category",
				Type:         "test_type",
				Subtype:      "test_subtype",
				Description:  "test_description",
			},
			mockBehavior: func(s *mock_service.MockProductsList, input domain.CreateProductInput) {
				s.EXPECT().Create(input).Return("34c8d3e6-b8d7-43dc-847e-5764c4114856", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":"34c8d3e6-b8d7-43dc-847e-5764c4114856"}`,
		},

		{
			name:                "Empty Fields",
			inputBody:           `{}`,
			mockBehavior:        func(s *mock_service.MockProductsList, input domain.CreateProductInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"Key: 'CreateProductInput.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'CreateProductInput.Price' Error:Field validation for 'Price' failed on the 'required' tag\nKey: 'CreateProductInput.Category' Error:Field validation for 'Category' failed on the 'required' tag\nKey: 'CreateProductInput.Type' Error:Field validation for 'Type' failed on the 'required' tag\nKey: 'CreateProductInput.Subtype' Error:Field validation for 'Subtype' failed on the 'required' tag"}`,
		},

		{
			name:      "Service Failure",
			inputBody: `{"title":"test_title","image":"test_image","price":70000,"sale":0,"sale_old_price":0,"category":"test_category","type":"test_type","subtype":"test_subtype","description":"test_description"}`,
			inputUser: domain.CreateProductInput{
				Title:        "test_title",
				Price:        70000,
				Sale:         0,
				SaleOldPrice: 0,
				Category:     "test_category",
				Type:         "test_type",
				Subtype:      "test_subtype",
				Description:  "test_description",
			},
			mockBehavior: func(s *mock_service.MockProductsList, input domain.CreateProductInput) {
				s.EXPECT().Create(input).Return("", errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			product := mock_service.NewMockProductsList(c)
			testCase.mockBehavior(product, testCase.inputUser)

			services := &service.Service{ProductsList: product}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/create-product", handler.createProduct)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/create-product", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getAllProducts(t *testing.T) {
	type mockBehavior func(s *mock_service.MockProductsList)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockProductsList) {
				s.EXPECT().GetAll().Return([]domain.ProductsList{
					{
						Id:           "453b4f0f-1f56-4c57-b43d-7b79792450a7",
						Title:        "Твидовый кардиган из хлопка",
						Image:        "w1.webp",
						Price:        749000,
						Sale:         0,
						SaleOldPrice: 0,
						Category:     "Женщинам",
						Type:         "Одежда",
						Subtype:      "Старые-коллекции",
						Description:  "",
					},
					{
						Id:           "b07221f8-4133-4688-b2d6-d677f41f5b74",
						Title:        "Объемный водоотталкивающий тренч",
						Image:        "w2.webp",
						Price:        499000,
						Sale:         50,
						SaleOldPrice: 999000,
						Category:     "Женщинам",
						Type:         "Одежда",
						Subtype:      "Старые-коллекции",
						Description:  "",
					},
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[{"id":"453b4f0f-1f56-4c57-b43d-7b79792450a7","title":"Твидовый кардиган из хлопка","image":"w1.webp","price":749000,"sale":0,"sale_old_price":0,"category":"Женщинам","type":"Одежда","subtype":"Старые-коллекции","description":"","created_at":"0001-01-01T00:00:00Z"},{"id":"b07221f8-4133-4688-b2d6-d677f41f5b74","title":"Объемный водоотталкивающий тренч","image":"w2.webp","price":499000,"sale":50,"sale_old_price":999000,"category":"Женщинам","type":"Одежда","subtype":"Старые-коллекции","description":"","created_at":"0001-01-01T00:00:00Z"}]}`,
		},

		{
			name: "Service Failure",
			mockBehavior: func(s *mock_service.MockProductsList) {
				s.EXPECT().GetAll().Return([]domain.ProductsList{}, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			product := mock_service.NewMockProductsList(c)
			testCase.mockBehavior(product)

			services := &service.Service{ProductsList: product}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.GET("/get-all-product", handler.getAllProducts)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/get-all-product", nil)

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getProductById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockProductsList, productId string)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockProductsList, productId string) {
				s.EXPECT().GetById(productId).Return(domain.ProductsList{
					Id:           "453b4f0f-1f56-4c57-b43d-7b79792450a7",
					Title:        "Твидовый кардиган из хлопка",
					Image:        "w1.webp",
					Price:        749000,
					Sale:         0,
					SaleOldPrice: 0,
					Category:     "Женщинам",
					Type:         "Одежда",
					Subtype:      "Старые-коллекции",
					Description:  "",
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":{"id":"453b4f0f-1f56-4c57-b43d-7b79792450a7","title":"Твидовый кардиган из хлопка","image":"w1.webp","price":749000,"sale":0,"sale_old_price":0,"category":"Женщинам","type":"Одежда","subtype":"Старые-коллекции","description":"","created_at":"0001-01-01T00:00:00Z"}}`,
		},

		{
			name: "Service Failure",
			mockBehavior: func(s *mock_service.MockProductsList, productId string) {
				s.EXPECT().GetById(productId).Return(domain.ProductsList{}, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			product := mock_service.NewMockProductsList(c)
			testCase.mockBehavior(product, "453b4f0f-1f56-4c57-b43d-7b79792450a7")

			services := &service.Service{ProductsList: product}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.GET("/get-product/:id", handler.getProductById)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/get-product/453b4f0f-1f56-4c57-b43d-7b79792450a7", nil)

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_updateProduct(t *testing.T) {
	type mockBehavior func(s *mock_service.MockProductsList, productId string, product domain.UpdateProductInput)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           domain.UpdateProductInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"title":"new_title","price":70000,"sale":0,"sale_old_price":0,"category":"new_category","type":"new_type","subtype":"new_subtype","description":"new_description"}`,
			inputUser: domain.UpdateProductInput{
				Title:        stringPointer("new_title"),
				Price:        uintPointer(70000),
				Sale:         uintPointer(0),
				SaleOldPrice: uintPointer(0),
				Category:     stringPointer("new_category"),
				Type:         stringPointer("new_type"),
				Subtype:      stringPointer("new_subtype"),
				Description:  stringPointer("new_description"),
			},
			mockBehavior: func(s *mock_service.MockProductsList, productId string, product domain.UpdateProductInput) {
				s.EXPECT().Update(productId, product).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"ok"}`,
		},

		{
			name:      "Service Failure",
			inputBody: `{"title":"new_title","price":70000,"sale":0,"sale_old_price":0,"category":"new_category","type":"new_type","subtype":"new_subtype","description":"new_description"}`,
			inputUser: domain.UpdateProductInput{
				Title:        stringPointer("new_title"),
				Price:        uintPointer(70000),
				Sale:         uintPointer(0),
				SaleOldPrice: uintPointer(0),
				Category:     stringPointer("new_category"),
				Type:         stringPointer("new_type"),
				Subtype:      stringPointer("new_subtype"),
				Description:  stringPointer("new_description"),
			},
			mockBehavior: func(s *mock_service.MockProductsList, productId string, product domain.UpdateProductInput) {
				s.EXPECT().Update(productId, product).Return(errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			product := mock_service.NewMockProductsList(c)
			testCase.mockBehavior(product, "453b4f0f-1f56-4c57-b43d-7b79792450a7", testCase.inputUser)

			services := &service.Service{ProductsList: product}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.PUT("/update-product/:id", handler.updateProduct)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/update-product/453b4f0f-1f56-4c57-b43d-7b79792450a7", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_deleteProduct(t *testing.T) {
	type mockBehavior func(s *mock_service.MockProductsList, productId string)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockProductsList, productId string) {
				s.EXPECT().Delete(productId).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"ok"}`,
		},

		{
			name: "Service Failure",
			mockBehavior: func(s *mock_service.MockProductsList, productId string) {
				s.EXPECT().Delete(productId).Return(errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			product := mock_service.NewMockProductsList(c)
			testCase.mockBehavior(product, "453b4f0f-1f56-4c57-b43d-7b79792450a7")

			services := &service.Service{ProductsList: product}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.DELETE("/delete-product/:id", handler.deleteProduct)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/delete-product/453b4f0f-1f56-4c57-b43d-7b79792450a7", nil)

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func uintPointer(b uint) *uint {
	return &b
}
