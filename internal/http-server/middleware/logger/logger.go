package logger

import (
"net/http"
"time"

"github.com/go-chi/chi/v5/middleware"
"golang.org/x/exp/slog"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			//создаем копию логгера добавляя подсказку, что это компонент мидлвейр логгера
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled") //эта строчка будет выведена при запуске приложения чтобы знать что такой хендлер у нас есть

		fn := func(w http.ResponseWriter, r *http.Request) {
			//эта часть будет находиться в цепочке хендлеров и выводиться при каждом входящем запросе
			entry := log.With( //будет выполнено до обработки запроса
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())), //присваивается реквестайди
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor) //врапер используется чтобы получить сведения об ответе которые мы используем далее в дефер

			t1 := time.Now()
			defer func() { //будет выполнено после обработки запроса
				entry.Info("request completed",
					slog.Int("status", ww.Status()), //статус запроса
					slog.Int("bytes", ww.BytesWritten()), //сколько байт записано
					slog.String("duration", time.Since(t1).String()), //время на обработку запроса
				)
			}()

			next.ServeHTTP(ww, r) //передаем управление следующему хендлеру в цепочке
		}

		return http.HandlerFunc(fn)
	}
}
