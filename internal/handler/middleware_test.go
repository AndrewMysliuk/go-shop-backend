package handler

import (
	"errors"
	"fmt"
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

func TestHandler_userIdentify(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
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
					CreatedAt: time.Now(),
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "user_id:34c8d3e6-b8d7-43dc-847e-5764c4114856, user_role:ADMIN",
		},

		{
			name:                 "No header",
			headerName:           "",
			mockBehavior:         func(s *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth header"}`,
		},

		{
			name:                 "Invalid Bearer",
			headerName:           "Authorization",
			headerValue:          "Bearr token",
			mockBehavior:         func(s *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid bearer"}`,
		},

		{
			name:                 "Empty Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			mockBehavior:         func(s *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty token"}`,
		},

		{
			name:        "Service Failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockUser, token string) {
				s.EXPECT().GetMe(token).Return(domain.User{}, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
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
			r.POST("/protected", handler.userIdentify, func(c *gin.Context) {
				userId, _ := c.Get(userCtx)
				role, _ := c.Get(userRole)

				c.String(200, fmt.Sprintf("user_id:%s, user_role:%s", userId, role))
			})

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
