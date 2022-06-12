package console

import (
	"errors"
	"os"
	"os/signal"

	"femalegeek/config"
	"femalegeek/db"
	"femalegeek/repository"
	"femalegeek/usecase"

	"femalegeek/service"

	"github.com/kumparan/cacher"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Long:  `This subcommand start the server`,
	Run:   run,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {
	signalChan := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		errCh <- errors.New("received an interrupt")
		db.StopTickerCh <- true
	}()

	db.InitializePostgresConn()
	keeper := cacher.NewKeeper()
	if !config.DisableCaching() {
		redisConn := db.NewRedisConnPool(config.RedisCacheHost())
		redisLockConn := db.NewRedisConnPool(config.RedisLockHost())

		keeper.SetConnectionPool(redisConn)
		keeper.SetLockConnectionPool(redisLockConn)
		keeper.SetDefaultTTL(config.CacheTTL())
	}
	keeper.SetDisableCaching(config.DisableCaching())

	userRepo := repository.NewUserRepository(db.DB, keeper)
	userUsecase := usecase.NewUserUsecase(userRepo)

	go func() {
		// Start HTTP server
		e := echo.New()
		e.Pre(middleware.AddTrailingSlash())
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

		svc := service.NewHTTPService(userUsecase)
		svc.Routes(e)

		errCh <- e.Start(":" + config.HTTPPort())
	}()

	log.Error(<-errCh)
}
