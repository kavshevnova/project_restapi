package slogdiscard

import (
	"log"
	"log/slog"
	"context"
)

func NewDiscardLogger() *slog.Logger {
	//собираем логгер и возвращаем его
	return slog.New(NewDiscardHandler())
}

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