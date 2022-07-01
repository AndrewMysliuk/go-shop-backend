package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/AndrewMislyuk/go-shop-backend/internal/service"
	mock_service "github.com/AndrewMislyuk/go-shop-backend/internal/service/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, user domain.UserSignUp)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           domain.UserSignUp
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test_Name","surname":"Test_Surname","email":"test@gmail.com","phone":"+4456780123","password":"1234QWER@","role":"ADMIN"}`,
			inputUser: domain.UserSignUp{
				Name:     "Test_Name",
				Surname:  "Test_Surname",
				Email:    "test@gmail.com",
				Phone:    "+4456780123",
				Password: "1234QWER@",
				Role:     "ADMIN",
			},
			mockBehavior: func(s *mock_service.MockUser, user domain.UserSignUp) {
				s.EXPECT().CreateUser(user).Return("34c8d3e6-b8d7-43dc-847e-5764c4114856", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":"34c8d3e6-b8d7-43dc-847e-5764c4114856"}`,
		},

		{
			name:                "Empty Fields",
			inputBody:           `{}`,
			mockBehavior:        func(s *mock_service.MockUser, user domain.UserSignUp) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"Key: 'UserSignUp.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'UserSignUp.Surname' Error:Field validation for 'Surname' failed on the 'required' tag\nKey: 'UserSignUp.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'UserSignUp.Phone' Error:Field validation for 'Phone' failed on the 'required' tag\nKey: 'UserSignUp.Password' Error:Field validation for 'Password' failed on the 'required' tag\nKey: 'UserSignUp.Role' Error:Field validation for 'Role' failed on the 'required' tag"}`,
		},

		{
			name:      "Service Failure",
			inputBody: `{"name":"Test_Name","surname":"Test_Surname","email":"test@gmail.com","phone":"+4456780123","password":"1234QWER@","role":"ADMIN"}`,
			inputUser: domain.UserSignUp{
				Name:     "Test_Name",
				Surname:  "Test_Surname",
				Email:    "test@gmail.com",
				Phone:    "+4456780123",
				Password: "1234QWER@",
				Role:     "ADMIN",
			},
			mockBehavior: func(s *mock_service.MockUser, user domain.UserSignUp) {
				s.EXPECT().CreateUser(user).Return("", errors.New("service failure"))
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

			auth := mock_service.NewMockUser(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{User: auth}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, user domain.UserSignIn)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           domain.UserSignIn
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"test@gmail.com","password":"1234QWER@"}`,
			inputUser: domain.UserSignIn{
				Email:    "test@gmail.com",
				Password: "1234QWER@",
			},
			mockBehavior: func(s *mock_service.MockUser, user domain.UserSignIn) {
				s.EXPECT().GenerateToken(user.Email, user.Password).Return("token", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"access_token":"token"}`,
		},

		{
			name:                "Empty Fields",
			inputBody:           `{}`,
			mockBehavior:        func(s *mock_service.MockUser, user domain.UserSignIn) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"Key: 'UserSignIn.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'UserSignIn.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},

		{
			name:      "Service Failure",
			inputBody: `{"email":"test@gmail.com","password":"1234QWER@"}`,
			inputUser: domain.UserSignIn{
				Email:    "test@gmail.com",
				Password: "1234QWER@",
			},
			mockBehavior: func(s *mock_service.MockUser, user domain.UserSignIn) {
				s.EXPECT().GenerateToken(user.Email, user.Password).Return("", errors.New("service failure"))
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

			auth := mock_service.NewMockUser(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{User: auth}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getMe(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, token string)

	testTable := []struct {
		name                string
		headerValue         string
		token               string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockUser, token string) {
				s.EXPECT().GetMe(token).Return(domain.User{
					Id:        "34c8d3e6-b8d7-43dc-847e-5764c4114856",
					Name:      "Test_Name",
					Surname:   "Test_Surname",
					Email:     "test@gmail.com",
					Phone:     "+4456781234",
					Role:      "ADMIN",
					Password:  "1234QWER@",
					CreatedAt: time.Date(2022, 01, 12, 13, 8, 21, 32963, time.Local),
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":"34c8d3e6-b8d7-43dc-847e-5764c4114856","name":"Test_Name","surname":"Test_Surname","email":"test@gmail.com","phone":"+4456781234","role":"ADMIN","password":"1234QWER@","created_at":"2022-01-12T13:08:21.000032963+02:00"}`,
		},

		{
			name:        "Service Failure",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockUser, token string) {
				s.EXPECT().GetMe(token).Return(domain.User{}, errors.New("service failure"))
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

			auth := mock_service.NewMockUser(c)
			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{User: auth}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/get-me", handler.getMe)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/get-me", nil)
			req.Header.Set("Authorization", testCase.headerValue)

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
