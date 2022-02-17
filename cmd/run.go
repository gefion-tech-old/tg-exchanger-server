package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/plugins"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/sqlstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	mine_plugin "github.com/gefion-tech/tg-exchanger-server/internal/plugins/mine"
	whitebit_plugin "github.com/gefion-tech/tg-exchanger-server/internal/plugins/whitebit"
)

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run server.",
		Long:  `...`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.Init()
			if _, err := toml.DecodeFile(
				fmt.Sprintf("config/config.%s.toml", viper.GetString("env")), cfg); err != nil {
				panic(err)
			}

			if err := runner(cfg); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().IntP("cpu", "c", 2, "Number of processor threads")
	cmd.Flags().StringP("env", "e", "local", "Launch environment")

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}

	return cmd
}

func runner(cfg *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	runtime.GOMAXPROCS(viper.GetInt("cpu"))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	postgres, err := db.InitPostgres(&cfg.Services.DB)
	if err != nil {
		panic(err)
	}
	defer postgres.Close()

	nsq, err := db.InitNSQ(&cfg.Services.NSQ)
	if err != nil {
		panic(err)
	}
	defer nsq.Stop()

	redisStore, closer, err := redisstore.InitAppRedisDictionaries(&cfg.Services.Redis)
	if err != nil {
		panic(err)
	}
	defer closer()

	sqlStore := sqlstore.Init(postgres)

	nsqStore := nsqstore.Init(nsq)

	plugins := plugins.InitAppPlugins(
		mine_plugin.InitMinePlugin(),
		whitebit_plugin.InitWhitebitPlugin(&cfg.Plugins),
	)

	logger := utils.InitLogger(sqlStore.AdminPanel().Logs())

	// lsnr := listener.InitListener(
	// 	sqlStore,
	// 	nsqStore,
	// 	plugins,
	// 	logger,
	// )

	srv := server.Init(
		sqlStore,
		nsqStore,
		redisStore,
		plugins,
		logger,
		cfg,
	).Create()

	// Запуск сервера
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.NewRecord(&models.LogRecord{
				Service: AppType.LogTypeServer,
				Module:  AppType.LogModuleServer,
				Info:    err.Error(),
			})
		}
	}()

	// Запуск слушателя транзакций
	// go func() {
	// 	if err := lsnr.Listen(ctx, &cfg.Listener); err != nil {
	// 		logger.NewRecord(&models.LogRecord{
	// 			Service: AppType.LogTypeServer,
	// 			Module:  AppType.LogModuleListener,
	// 			Info:    err.Error(),
	// 		})
	// 	}
	// }()

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.NewRecord(&models.LogRecord{
			Service: AppType.LogTypeServer,
			Module:  AppType.LogModuleServer,
			Info:    "Server forced to shutdown: " + err.Error(),
		})
	}

	return nil
}
