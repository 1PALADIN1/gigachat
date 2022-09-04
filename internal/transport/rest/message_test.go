package rest

import (
	"errors"
	"fmt"
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

func TestHandler_getAllChatMessages(t *testing.T) {
	type mockBehavior func(m *mock_service.MockMessage, a *mock_service.MockAuthorization, chatId int)

	testTable := []struct {
		name                 string
		inputChatId          int
		headerName           string
		headerValue          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputChatId: 1,
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(m *mock_service.MockMessage, a *mock_service.MockAuthorization, chatId int) {
				a.EXPECT().ParseToken("token").Return(1, nil)
				m.EXPECT().GetAllMessages(chatId).Return([]entity.Message{
					{
						Id:       1,
						SendTime: time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
						Text:     "Hello!",
						UserId:   1,
						Username: "test_user",
						ChatId:   1,
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"send_time":"2022-08-31 00:00:00","text":"Hello!","chat_id":1,"user_id":1,"username":"test_user"}]`,
		},
		{
			name:        "OK Empty",
			inputChatId: 1,
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(m *mock_service.MockMessage, a *mock_service.MockAuthorization, chatId int) {
				a.EXPECT().ParseToken("token").Return(1, nil)
				m.EXPECT().GetAllMessages(chatId).Return([]entity.Message{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[]`,
		},
		{
			name:        "Invalid Auth Header",
			inputChatId: 1,
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(m *mock_service.MockMessage, a *mock_service.MockAuthorization, chatId int) {
				a.EXPECT().ParseToken("token").Return(0, errors.New("invalid auth header"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:        "Service Failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(m *mock_service.MockMessage, a *mock_service.MockAuthorization, chatId int) {
				a.EXPECT().ParseToken("token").Return(1, nil)
				m.EXPECT().GetAllMessages(chatId).Return(nil, errors.New("service internal error"))
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
			messageService := mock_service.NewMockMessage(c)
			testCase.mockBehavior(messageService, authService, testCase.inputChatId)

			services := &service.Service{
				Authorization: authService,
				Message:       messageService,
			}

			handler := NewHandler(services)

			// Setup Server
			r := mux.NewRouter()
			r.HandleFunc("/api/chat/{id:[0-9]+}/message", handler.getAllChatMessages).Methods(http.MethodGet)

			// Perform Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/chat/%d/message", testCase.inputChatId), nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)
			r.ServeHTTP(w, req)

			// Validate
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
