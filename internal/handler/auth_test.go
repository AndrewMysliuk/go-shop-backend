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
