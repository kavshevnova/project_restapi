package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	//часть ответа которая будет повторяться для всех хендлеров
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK = "ok"
	StatusError = "error"
)

func OK() Response {
	return Response{Status: StatusOK}
}

func Error(msg string) Response {
	return Response{Status: StatusError, Error: msg}
}

//принимаем список ошибок валидатора и возвращаем ответ
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	//перебираем все ошибки которые мы получили и формируем ответ клиенту
	for _, err := range errs {
		switch err.ActualTag() {//смотрим что за тег у validator.ValidationErrors []FieldError
		case "required": //такое поле было обязательным
			errMsgs = append(errMsgs, fmt.Sprintf("failed %s is a required field", err.Field()))
		case "url": //это поле не валидно юрл
			errMsgs = append(errMsgs, fmt.Sprintf("failed %s is not a valid URL field", err.Field()))
			default: //поле не валидное
				errMsgs = append(errMsgs, fmt.Sprintf("failed %s is not valid", err.Field()))
		}
	}

	return Response{Status: StatusError, Error: strings.Join(errMsgs, ", ")}
}