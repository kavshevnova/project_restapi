package redirect_test

import (
	"github.com/go-chi/chi/v5"
	"github.com/kavshevova/project_restapi/internal/http-server/handlers/redirect/mocks"
	"github.com/kavshevova/project_restapi/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/kavshevova/project_restapi/internal/http-server/handlers/redirect"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "все нормально",
			alias: "какой-то алиас",
			url:   "http://example.com/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URLGetterMock := mocks.NewURLGetter(t)
			if tt.mockError != nil || tt.respError == "" {
				URLGetterMock.On("Get", tt.alias).
					Return(tt.url, tt.mockError).
					Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), URLGetterMock)

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectToUrl, err := api.GetRedirectURL(ts.URL + "/" + tt.alias)
			require.NoError(t, err)

			assert.Equal(t, tt.url, redirectToUrl)

		})
	}
}
