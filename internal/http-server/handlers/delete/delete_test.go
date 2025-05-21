package delete_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	delete2 "github.com/kavshevova/project_restapi/internal/http-server/handlers/delete"
	"github.com/kavshevova/project_restapi/internal/http-server/handlers/delete/mocks"
	resp "github.com/kavshevova/project_restapi/internal/lib/api/response"
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
		wantCode int
		want resp.Response
		mockError error
		setupMock bool
	}{
		{
			name:      "пустой алиас",
			alias:     "",
			wantCode:  http.StatusBadRequest,
			want:      resp.Error("invalid request"),
			setupMock: false,
		},
		{
			name:   "все нормально",
			alias:  "какой-то алиас",
			wantCode: http.StatusOK,
			want: resp.OK(),
			setupMock: true,
		},
		{
			name:      "внутренняя ошибка",
			alias:     "какой-то алиас",
			wantCode:  http.StatusInternalServerError,
			want:      resp.Error("internal server error"),
			mockError: errors.New("internal server error"),
			setupMock: true,
		},
	}
	for _, testCase := range tests {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			URLDeleterMock := mocks.NewURLDeleter(t)

			if testCase.setupMock {
				URLDeleterMock.On("DeleteURL", testCase.alias).
					Return(testCase.mockError).
					Once()
			}

			handler := delete2.New(slogdiscard.NewDiscardLogger(), URLDeleterMock)
			req, err := http.NewRequest(http.MethodDelete, "/"+testCase.alias, nil)
			require.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", testCase.alias)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

				var resp resp.Response
				err = json.NewDecoder(rr.Body).Decode(&resp)
				require.NoError(t, err, "ошибка декодирования JSON")
				require.Equal(t, testCase.want, resp)

			if testCase.setupMock {
				URLDeleterMock.AssertExpectations(t)
			}
		})
	}
}