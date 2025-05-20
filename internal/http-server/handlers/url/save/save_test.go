package save_test
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kavshevova/project_restapi/internal/http-server/handlers/url/save"
	"github.com/kavshevova/project_restapi/internal/http-server/handlers/url/save/mocks"
	"github.com/kavshevova/project_restapi/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)
func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "success",
			url:   "http://example.com/",
			alias: "example",
		},
		{
			name:  "empty alias",
			url:   "http://example.com/",
			alias: "",
		},
		{
			name:      "empty url",
			url:       "",
			alias:     "example",
			respError: "url cannot be empty",
		},
		{
			name:      "invalid url",
			url:       "example.com",
			alias:     "example",
			respError: "url contains invalid characters",
		},
		{
			name:      "SaveURL Error",
			url:       "http://example.com/",
			alias:     "example",
			respError: "failed to save url",
			mockError: errors.New("unexpected error"),
		},
	}
	for _, testCase := range tests {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
 t.Parallel()

 URLSaverMock := mocks.NewURLSaver(t) //создаем объект мока который мы сгенерировали

 if testCase.mockError != nil || testCase.respError == "" {
	 URLSaverMock.On("SaveURL", testCase.url, mock.AnythingOfType("string")).
		 Return(int64(1), testCase.mockError).
		 Once()
 }

 handler := save.New(slogdiscard.NewDiscardLogger(), URLSaverMock)

 //пример запроса, который отправим (он будет в виде джсона и у него будут параметры которые перечислены в тесткейсе)
 input := fmt.Sprintf(`{"url":"%s", "alias":"%s"}`, testCase.url, testCase.alias)
 //создаем новый запрос
 req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
 require.NoError(t, err) //если есть ошибка то тест фейлится
 //assert.NoError() выводит ошибку но тест не фейлится а идет дальше

//выполняется запрос
 rr := httptest.NewRecorder() //создаем респонсрекордер чтобы записать туда ответ нашего хендлера
 handler.ServeHTTP(rr, req) //запускаем запрос


 require.Equal(t, rr.Code, http.StatusOK) //смотрим что было записано в рекордер
 body := rr.Body.String() //смотрим что было записано в тело
 var resp save.Response //создаем объект в который запишем ответ
 require.NoError(t, json.Unmarshal([]byte(body), &resp)) //анмаршалим
 require.Equal(t, testCase.respError, resp.Error) //смотрим что ошибка которую вернул хендлер совпадает с ошибкой в тесткейсе
		})
	}
}