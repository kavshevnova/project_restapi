package delete_test

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/kavshevova/project_restapi/internal/http-server/handlers/delete"
	"github.com/kavshevova/project_restapi/internal/http-server/handlers/delete/mocks"
	"github.com/kavshevova/project_restapi/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		respError string
		respOk    string
		mockError error
	}{
		{
			name:      "пустой алиас",
			alias:     "",
			respError: "пустой алиас",
		},
		{
			name:   "все нормально",
			alias:  "какой-то алиас",
			respOk: "ок",
		},
		{
			name:      "внутренняя ошибка",
			alias:     "какой-то алиас",
			mockError: errors.New("внутренняя ошибка"),
		},
	}
	for _, testCase := range tests {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			URLDeleterMock := mocks.NewURLDeleter(t)

			if testCase.mockError != nil || testCase.respError == "" {
				URLDeleterMock.On("Delete", testCase.alias).
					Return(testCase.mockError).
					Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), URLDeleterMock)

			req, err := http.NewRequest(http.MethodDelete, "/delete", nil)
			require.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", testCase.alias)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, testCase.respOk, rr.Body.String())

		})
	}
}