package main

import (
	"github.com/joho/godotenv"
	"github.com/kavshevova/project_restapi/internal/config"
	"log"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("Starting server", slog.String("env", cfg.Env))
	log.Debug("debug logging enabled")
	//TODO init storage: sqlite3
	//TODO: init router: chi, "chi render"
	//TODO: run server
}

func setupLogger (env string) *slog.Logger {
	//почему логгер должен зависеть от параметра енв:
	//локально мы хотим видеть текстовые логи
	//в окружении дев или прод (на сервере) мы хотим видеть json. Причем на дев - логи уровня дебаг, а на проде не ниже уровня инфо.
	var log *slog.Logger
	//объявили логгер и в зависимости от переменной енв мы будем его создавать
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
		case envDev:
			log = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
			case envProd:
				log = slog.New(
					slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
					)
	}
	return log
}