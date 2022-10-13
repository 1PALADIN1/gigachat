package rest

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/1PALADIN1/gigachat_server/internal/entity"
	"github.com/1PALADIN1/gigachat_server/internal/service"
	mock_service "github.com/1PALADIN1/gigachat_server/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_createChat(t *testing.T) {
	type mockBehavior func(c *mock_service.MockChat, a *mock_service.MockAuthorization, chat entity.Chat)

	testTable := []struct {
		name                 string
		inputBody            string
		inputChat            entity.Chat
		headerName           string
		headerValue          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"title":"test_chat","user_ids":[1,2]}`,
			inputChat: entity.Chat{
				Title:   "test_chat",
				UserIds: []int{1, 2},
			},
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization, chat entity.Chat) {
				c.EXPECT().GetOrCreateChat(chat).Return(1, nil)
				a.EXPECT().ParseToken("token").Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "OK with Description",
			inputBody: `{"title":"test_chat","description":"test_description","user_ids":[1,2]}`,
			inputChat: entity.Chat{
				Title:       "test_chat",
				Description: "test_description",
				UserIds:     []int{1, 2},
			},
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization, chat entity.Chat) {
				c.EXPECT().GetOrCreateChat(chat).Return(1, nil)
				a.EXPECT().ParseToken("token").Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "OK Group Users",
			inputBody: `{"title":"test_chat","user_ids":[1,2,3]}`,
			inputChat: entity.Chat{
				Title:   "test_chat",
				UserIds: []int{1, 2, 3},
			},
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization, chat entity.Chat) {
				c.EXPECT().GetOrCreateChat(chat).Return(1, nil)
				a.EXPECT().ParseToken("token").Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:        "Bad JSON",
			inputBody:   "",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization, chat entity.Chat) {
				a.EXPECT().ParseToken("token").Return(1, nil)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid request body"}`,
		},
		{
			name:        "Empty Title",
			inputBody:   `{"description":"test_description","user_ids":[1,2]}`,
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization, chat entity.Chat) {
				a.EXPECT().ParseToken("token").Return(1, nil)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"chat title is not set"}`,
		},
		{
			name:        "Invalid Users Amount",
			inputBody:   `{"title":"test_chat","user_ids":[1]}`,
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization, chat entity.Chat) {
				a.EXPECT().ParseToken("token").Return(1, nil)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"need at least 2 users to create chat"}`,
		},
		{
			name:        "Empty Users",
			inputBody:   `{"title":"test_chat"}`,
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization, chat entity.Chat) {
				a.EXPECT().ParseToken("token").Return(1, nil)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"need at least 2 users to create chat"}`,
		},
		{
			name:        "Duplicate Users",
			inputBody:   `{"title":"test_chat","user_ids":[1,2,1]}`,
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization, chat entity.Chat) {
				a.EXPECT().ParseToken("token").Return(1, nil)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"duplicate user id=1 in request"}`,
		},
		{
			name:        "Duplicate Users #2",
			inputBody:   `{"title":"test_chat","user_ids":[2,2]}`,
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization, chat entity.Chat) {
				a.EXPECT().ParseToken("token").Return(1, nil)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"duplicate user id=2 in request"}`,
		},
		{
			name:        "Invalid Auth Header",
			inputBody:   `{"title":"test_chat","user_ids":[1,2]}`,
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization, chat entity.Chat) {
				a.EXPECT().ParseToken("token").Return(0, errors.New("invalid auth header"))
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

			chatService := mock_service.NewMockChat(c)
			authService := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(chatService, authService, testCase.inputChat)

			services := &service.Service{
				Authorization: authService,
				Chat:          chatService,
			}

			handler := NewHandler(services)

			// Setup Server
			r := mux.NewRouter()
			r.HandleFunc("/api/chat", handler.createChat).Methods(http.MethodPost)
			r.Use(handler.validateAuthHeader)

			//Perform Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set(testCase.headerName, testCase.headerValue)
			r.ServeHTTP(w, req)

			//Validate
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getAllChats(t *testing.T) {
	type mockBehavior func(c *mock_service.MockChat, a *mock_service.MockAuthorization)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization) {
				a.EXPECT().ParseToken("token").Return(1, nil)
				c.EXPECT().GetAllChats(1).Return([]entity.ChatResponse{
					{
						Chat: entity.Chat{
							Id:          1,
							Title:       "test_chat",
							Description: "test_description",
						},
						LastMessage:         "Hello!",
						LastMessageUserId:   "1",
						LastMessageUsername: "test_user",
						LastMessageTime:     time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC).String(),
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"id":1,"title":"test_chat","description":"test_description","last_message":"Hello!","last_message_user_id":"1","last_message_username":"test_user","last_message_time":"2022-08-31 00:00:00 +0000 UTC"}]`,
		},
		{
			name:        "Invalid Auth Header",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization) {
				a.EXPECT().ParseToken("token").Return(0, errors.New("invalid auth header"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:        "Service Failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(c *mock_service.MockChat, a *mock_service.MockAuthorization) {
				a.EXPECT().ParseToken("token").Return(1, nil)
				c.EXPECT().GetAllChats(1).Return(nil, errors.New("service internal error"))
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

			chatService := mock_service.NewMockChat(c)
			authService := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(chatService, authService)

			services := &service.Service{
				Authorization: authService,
				Chat:          chatService,
			}

			handler := NewHandler(services)

			// Setup Server
			r := mux.NewRouter()
			r.HandleFunc("/api/chat", handler.getAllChats).Methods(http.MethodGet)
			r.Use(handler.validateAuthHeader)

			// Perform Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/chat", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)
			r.ServeHTTP(w, req)

			// Validate
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
