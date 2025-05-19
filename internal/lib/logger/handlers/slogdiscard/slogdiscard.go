package slogdiscard

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	//собираем логгер и возвращаем его
	return slog.New(NewDiscardHandler())
}

//1. Метод Enabled()
//Enabled(context.Context, Level) bool
//Назначение: Проверяет, нужно ли логировать сообщение данного уровня.
//Параметры:
//context.Context — контекст запроса (может содержать метаданные)
//Level — уровень логирования (Debug, Info, Warn, Error)
//Возвращает: true, если сообщение нужно обработать.
//2. Метод Handle()
//Handle(context.Context, Record) error
//Назначение: Обрабатывает запись лога.
//Параметры:
//context.Context — контекст
//Record — структура с данными лога.
//Возвращает: Ошибку, если запись не удалась.
//3. Метод WithAttrs()
//WithAttrs(attrs []Attr) Handler
//Назначение: Создает новый обработчик с добавленными атрибутами.
//Параметры:
//attrs — массив атрибутов вида slog.String("key", "value")
//Возвращает: Новый обработчик (старый остается неизменным).
//4. Метод WithGroup()
//WithGroup(name string) Handler
//Назначение: Группирует последующие атрибуты под заданным именем.
//Параметры:
//name — название группы
//Возвращает: Новый обработчик с поддержкой группировки.


type DiscardHandler struct {} //вместо того чтобы хендлить сообщения он будет их игнорировать

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, r slog.Record) error {
	//просто игнорируем запись журнала
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	//возвращает тот же обработчик так как нет атрибутов для сохранения
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	//возвращает тот же обработчикБ так как нет группы для сохранения
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	//всегда возвращает false так как запись журнала игнорируется
	return false
}