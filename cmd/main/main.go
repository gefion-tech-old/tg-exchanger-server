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
	prod bool
	cpu  int
)

func init() {
	flag.BoolVar(&prod, "prod", false, "Strat on production server.")
	flag.IntVar(&cpu, "cpu", 2, "Number of processor threads")
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(cpu)
	ctx := context.Background()

	// Инициализирую конфигурацию
	var cnf *config.Config
	if prod {
		cnf = config.Init()
		if _, err := toml.DecodeFile("config/config.prod.toml", cnf); err != nil {
			panic(err)
		}

	} else {
		cnf = config.Init()
		if _, err := toml.DecodeFile("config/config.local.toml", cnf); err != nil {
			panic(err)
		}
	}

	// Создаю подключение к Postgres
	postgres, err := db.InitPostgres(&cnf.DB)
	if err != nil {
		panic(err)
	}
	defer postgres.Close()

	// Инициализация соединения с NSQ
	nsq, err := db.InitNSQ(&cnf.NSQ)
	if err != nil {
		panic(err)
	}

	// Инициализация модуля приложения
	application := app.Init(postgres, nsq, cnf)
	if err := application.Start(ctx); err != nil {
		fmt.Println(err)
	}

}
