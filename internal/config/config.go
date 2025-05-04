package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

//здесь будет парсинг конфига

type Config struct {
	//объект конфига будет полностью соответствовать ямл файлу
	/*как записано в ямл*/ /*значение по умолчанию*/ /*env-required: "true" - можно сделать и если мы забыли указать енв то на проде ничего не запустится, так мы точно не запустим случайно приложение в режиме локал у себя на проде*/
	Env string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"localhost:8080"`
	Timeout time.Duration    `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	//приставка маст в функциях используется когда вместо возврата ошибки функция будет паниковать
	//это функция, которая прочитает файл с конфига и создаст и заполнит объект конфиг
	configPath := os.Getenv("CONFIG_PATH")
	//считываем файл с конфигом из переменной окружения CONFIG_PATH
	if configPath == "" {
		//если там не находим то роняем приложение
		log.Fatal("CONFIG_PATH is not set")
	}
	 if _, err := os.Stat(configPath); os.IsNotExist(err) {
		 //проверяем существует ли такой файл
		 log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	 }

	 var cfg Config

	 if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		 log.Fatalf("Error reading config: %s", err)
	 }
 return &cfg

}