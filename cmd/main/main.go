package main

import (
	"flag"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

var (
	configPath string
	proc       int
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/config.local.toml", "Path to config file")
	flag.IntVar(&proc, "proc", 2, "Number of processor threads")
}

func main() {
	runtime.GOMAXPROCS(proc)

	// Инициализирую конфигурацию
	config := config.Init()
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		panic(err)
	}

	// Создаю подключение к Postgres
	postgres, err := db.InitPostgres(&config.DB)
	if err != nil {
		panic(err)
	}
	defer postgres.Close()

	// Создаю подключение к Redis
	redis, err := db.InitRedis(&config.Redis)
	if err != nil {
		panic(err)
	}
	defer redis.Close()

}
