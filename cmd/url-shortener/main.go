package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kavshevova/project_restapi/internal/config"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	cfg := config.MustLoad()

	fmt.Println(cfg)
	//TODO: init logger: slog
	//TODO init storage: sqlite3
	//TODO: init router: chi, "chi render"
	//TODO: run server
}
