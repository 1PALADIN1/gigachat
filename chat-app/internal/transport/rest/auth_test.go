package rest

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/service"
	mock_service "github.com/1PALADIN1/gigachat_server/internal/service/mocks"
	"github.com/1PALADIN1/gigachat_server/internal/transport/helper"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_singUpUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user entity.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            entity.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test_username","password":"test_password"}`,
			inputUser: entity.User{
				Username: "test_username",
				Password: "test_password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user entity.User) {
				s.EXPECT().SignUpUser(user).Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Bad JSON",
			inputBody:            "",
			mockBehavior:         func(s *mock_service.MockAuthorization, user entity.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid request body"}`,
		},
		{
			name:                 "Empty Username",
			inputBody:            `{"password":"test_password"}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, user entity.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"username is not set"}`,
		},
		{
			name:                 "Empty Password",
			inputBody:            `{"username":"test_username"}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, user entity.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"password is not set"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"username":"test_username","password":"test_password"}`,
			inputUser: entity.User{
				Username: "test_username",
				Password: "test_password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user entity.User) {
				s.EXPECT().SignUpUser(user).Return(0, errors.New("service internal error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"service internal error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{
				Authorization: auth,
			}

			handler := NewHandler(services)

			// Setup Server
			r := mux.NewRouter()
			r.HandleFunc("/api/auth/sign-up", handler.singUpUser).Methods(http.MethodPost)

			// Perform Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBufferString(testCase.inputBody))
			r.ServeHTTP(w, req)

			// Validate
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_signInUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user entity.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            entity.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test_username","password":"test_password"}`,
			inputUser: entity.User{
				Username: "test_username",
				Password: "test_password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user entity.User) {
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("token", 1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"access_token":"token"}`,
		},
		{
			name:                 "Bad JSON",
			inputBody:            "",
			mockBehavior:         func(s *mock_service.MockAuthorization, user entity.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid request body"}`,
		},
		{
			name:                 "Empty Username",
			inputBody:            `{"password":"test_password"}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, user entity.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"username is not set"}`,
		},
		{
			name:                 "Empty Password",
			inputBody:            `{"username":"test_username"}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, user entity.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"password is not set"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"username":"test_username","password":"test_password"}`,
			inputUser: entity.User{
				Username: "test_username",
				Password: "test_password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user entity.User) {
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("", 0, errors.New("service internal error"))
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

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{
				Authorization: auth,
			}

			handler := NewHandler(services)

			// Setup Server
			r := mux.NewRouter()
			r.HandleFunc("/api/auth/sign-in", handler.signInUser).Methods(http.MethodPost)

			// Perform Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewBufferString(testCase.inputBody))
			r.ServeHTTP(w, req)

			// Validate
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestValidateAuthHeader(t *testing.T) {
	type mockBehavior func(a *mock_service.MockAuthorization, token string)

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
			mockBehavior: func(a *mock_service.MockAuthorization, token string) {
				a.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Empty Header",
			headerName:           "",
			headerValue:          "",
			mockBehavior:         func(a *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                 "Empty Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			mockBehavior:         func(a *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                 "Empty Header Value",
			headerName:           "Authorization",
			headerValue:          "",
			mockBehavior:         func(a *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                 "Not Enough Header Parts",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			mockBehavior:         func(a *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                 "Not Enough Header Parts #2",
			headerName:           "Authorization",
			headerValue:          "Bearer   ",
			mockBehavior:         func(a *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                 "Unexpected Parts Amount",
			headerName:           "Authorization",
			headerValue:          "Bearer token token",
			mockBehavior:         func(a *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:        "Invalid Bearer",
			headerName:  "Authorization",
			headerValue: "Bearr token",
			mockBehavior: func(a *mock_service.MockAuthorization, token string) {
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			authService := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(authService, testCase.token)

			// Setup Server
			r := mux.NewRouter()
			r.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
				userId, ok := helper.ValidateAuthHeader(w, r, authService)
				if !ok {
					return
				}

				helper.SendResponse(w, http.StatusOK, map[string]interface{}{
					"id": userId,
				})
			}).Methods(http.MethodGet)

			// Perform Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/validate", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)
			r.ServeHTTP(w, req)

			// Validate
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
