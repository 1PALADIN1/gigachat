package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1PALADIN1/gigachat_server/internal/service"
	mock_service "github.com/1PALADIN1/gigachat_server/internal/service/mocks"
	"github.com/1PALADIN1/gigachat_server/internal/transport/helper"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

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

			services := &service.Service{
				Authorization: authService,
			}

			handler := NewHandler(services)

			// Setup Server
			r := mux.NewRouter()
			r.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
				userId, err := handler.parseHeader(r)
				if err != nil {
					helper.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
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
