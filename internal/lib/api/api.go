package api

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	//редирект будет происходить именно с кодом 302, если там будет друггой код то функция вернет ошибку
	ErrInvalidStatusCode = errors.New("invalid status code")
)

func GetRedirect(url string) (string, error) {
	const op = "api.GetRedirect"

	//создаем клиент кастомный чттп клиенту который будет проверять редиректыот и он ключает автоматическое следование редиректам
	client := &http.Client{
		//req — новый запрос, который будет выполнен
		//via — список уже выполненных запросов в цепочке
		//Возвращает:
		//nil — продолжить редирект
		//http.ErrUseLastResponse — остановиться
		//Другую error — прервать с ошибкой
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}, //При получении кода 301/302/307 клиент не будет автоматически переходить по новому URL
		//Вместо этого он:
		//Остановится на первом же редиректе
		//Вернёт исходный ответ с кодом перенаправления
		//Позволит вам вручную обработать заголовок Location
	}
	//делаем запрос на переданную функцию юрл
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	//проверяем что статус ответа совпадает с тем который мы ожидаем
	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("#{op}: #{ErrInvalidStatusCode}: #{resp.StatusCode}")
	}
	defer func() {_ = resp.Body.Close() }() //закрываем тело ответа

	//возвращаем тот юрл на который происходит редирект
	return resp.Header.Get("Location"), nil
}