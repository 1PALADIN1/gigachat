package rest

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/service"
	mock_service "github.com/1PALADIN1/gigachat_server/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_findUserByName(t *testing.T) {
	type mockBehavior func(u *mock_service.MockUser, a *mock_service.MockAuthorization, filter string, userId int)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		inputFilter          string
		inputUserId          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			inputFilter: "user",
			inputUserId: 1,
			mockBehavior: func(u *mock_service.MockUser, a *mock_service.MockAuthorization, filter string, userId int) {
				a.EXPECT().ParseToken("token").Return(1, nil)
				u.EXPECT().FindUserByName(filter, userId).Return([]entity.User{
					{
						Id:       2,
						Username: "user_2",
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"id":2,"username":"user_2"}]`,
		},
		{
			name:        "OK Not Found",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			inputFilter: "user",
			inputUserId: 1,
			mockBehavior: func(u *mock_service.MockUser, a *mock_service.MockAuthorization, filter string, userId int) {
				a.EXPECT().ParseToken("token").Return(userId, nil)
				u.EXPECT().FindUserByName(filter, userId).Return([]entity.User{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[]`,
		},
		{
			name:        "Invalid Auth Header",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			inputFilter: "user",
			inputUserId: 1,
			mockBehavior: func(u *mock_service.MockUser, a *mock_service.MockAuthorization, filter string, userId int) {
				a.EXPECT().ParseToken("token").Return(0, errors.New("invalid auth header"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:        "Service Failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			inputFilter: "user",
			inputUserId: 1,
			mockBehavior: func(u *mock_service.MockUser, a *mock_service.MockAuthorization, filter string, userId int) {
				a.EXPECT().ParseToken("token").Return(userId, nil)
				u.EXPECT().FindUserByName(filter, userId).Return(nil, errors.New("service internal error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"service internal error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			authService := mock_service.NewMockAuthorization(c)
			userService := mock_service.NewMockUser(c)
			testCase.mockBehavior(userService, authService, testCase.inputFilter, testCase.inputUserId)

			services := &service.Service{
				Authorization: authService,
				User:          userService,
			}

			handler := NewHandler(services)

			// Setup Server
			r := mux.NewRouter()
			r.HandleFunc("/api/user/{user}", handler.findUserByName).Methods(http.MethodGet)

			// Perform Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/user/%s", testCase.inputFilter), nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)
			r.ServeHTTP(w, req)

			// Validate
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
