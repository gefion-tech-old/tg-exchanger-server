package main

import (
	"context"
	"flag"
	"fmt"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/gefion-tech/tg-exchanger-server/internal/app"
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
	ctx := context.Background()

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

	// Инициализация соединения с NSQ
	nsq, err := db.InitNSQ(&config.NSQ)
	if err != nil {
		panic(err)
	}

	// Инициализация модуля приложения
	application := app.Init(postgres, nsq, config)
	if err := application.Start(ctx); err != nil {
		fmt.Println(err)
	}

}
