package tests

import (
	"github.com/brianvoe/gofakeit/v6" //библиотека которая позволяет генерировать случайные данные (имена, фамилии, mail, номера телефонов и тд)
	"github.com/gavv/httpexpect/v2"   //некий фреймворк который нужен чтобы тестировать http-сервисы
	"github.com/kavshevova/project_restapi/internal/http-server/handlers/url/save"
	"github.com/kavshevova/project_restapi/internal/lib/api"
	"github.com/kavshevova/project_restapi/internal/lib/random"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"path"
	"testing"

)

const host = "localhost:8082"

func TestURLShortener_HappyPath(t *testing.T) {
	//формируем базовый урл к которому будет обращаться наш клиент
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	//создаем специальный клиент из библиотеки httpexpect, с помощью которого будут отправляться запросы
	client := httpexpect.Default(t, u.String())

	client.POST("/url"). //формируем пост-запрос
		WithJSON(save.Request{ //передаем сюда объект запроса который дальше будет маршалиться в джсон
			URL: gofakeit.URL(), //генерируем случайный url
			Alias: random.NewRandomString(10), //генерируем случайный алиас функцией из моей библиотеки
	}).
		WithBasicAuth("myuser", "mypass"). //чтобы прошла авторизация
		Expect(). //что ожидать от ответа
		Status(http.StatusOK). //ожидаем статус 200
		JSON(). //формируем его в джсон
		Object(). //из джсона получаем объект
		ContainsKey("alias") //проверяем что этот ответ содержит параметр алиас
}

func TestURLShortener_SaveRedirect(t *testing.T) {
	testCase := []struct {
		name     string
		url      string
		alias    string
		error     string
	}{
		{
			name:  "valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word(),
		},
		{
			name:  "invalid URL",
			url:   "invalid-url",
			alias: gofakeit.Word(),
			error: "failed URL is not a valid URL field",
		},
		{
			name: "empty alias",
			url:  gofakeit.URL(),
			alias: "",
		},
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}
			client := httpexpect.Default(t, u.String())

			//resp - объект который мы получаем в ответ пост-запроса
			resp := client.POST("/url").
				WithJSON(save.Request{
					URL:   tt.url,
					Alias: tt.alias,
				}).
				WithBasicAuth("myuser", "mypass").
				Expect().
				Status(http.StatusOK).
				JSON().
				Object()

			//смотрим не было ли в тесткейсе указано что мы ожидаем ошибку
			//если мы ожидаем ошибку то мы ожидаем что в ответе не будет элиаса и смотрим что в ответе есть параметр error и он равен той ошибке которую мы указали в тесткейсе
			if tt.error != "" {
				resp.NotContainsKey("alias")
				resp.Value("error").String().IsEqual(tt.error)
				return //потому что дальше проверять нет смысла
			}
			resp.ContainsKey("alias")

			//получаем  алиас который у нас был указан в тесткейсе
			alias := tt.alias

			//если алиас не пустой то мы проверяемм что в теле ответа именно этот алиас и содержится
			if alias != "" {
				resp.Value("alias").String().IsEqual(tt.alias)
			} else {
				//если нет то то наш хендлер сгенерирует случайный алиас и его вернет
				//в это случае мы смотрим что пришло и сохраняеммв свою переменную алиас чтобы дальше по нему обратиться к сервису и получить редирект
				resp.Value("alias").String().NotEmpty()
				alias = resp.Value("alias").String().Raw()
			}

			//проверяем редирект
			testRedirect(t, alias, tt.url)

			reqDel := client.DELETE("/"+path.Join("url", alias)).
				WithBasicAuth("myuser", "mypass").
				Expect().Status(http.StatusOK).
				JSON().Object()
			reqDel.Value("status").String().IsEqual("ok")

			testRedirectNotFound(t, alias)

		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string)  {
	//формируем юрл
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	//вызываем функцию гет редирект и смотрим куда она будет редиректить
	redirecttourl, err := api.GetRedirect(u.String())
	//проверяем что не произошла какая-то ошибка
	require.NoError(t, err)
	//проверяем что полученный редирект совпадает с тем что мы получили в аргументе
	require.Equal(t, urlToRedirect, redirecttourl)

}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}
	_, err := api.GetRedirect(u.String())
	require.Equal(t, err, api.ErrInvalidStatusCode)


}