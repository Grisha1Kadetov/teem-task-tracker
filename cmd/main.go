package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Grisha1Kadetov/TeemTaskTrackerService/internal/application"
	"github.com/Grisha1Kadetov/TeemTaskTrackerService/internal/config"
	"github.com/Grisha1Kadetov/TeemTaskTrackerService/internal/pkg/log"
	//"github.com/Grisha1Kadetov/TeemTaskTrackerService/internal/pkg/db"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := log.NewZapLogger()
	defer logger.Close()

	conf, err := config.LoadConfig(logger)
	if err != nil {
		panic(err)
	}

	server := prepareServer(ctx, conf, logger)

	go func() {
		logger.Info("starting server", log.Pair("port", conf.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			stop()
			logger.Panic("failed to start server", log.Err(err))
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Panic("failed to shutdown server", log.Err(err))
	}
}

func prepareServer(ctx context.Context, cnf *config.Config, logger log.Logger) *http.Server {
	sqlDB := prepareDatabase(cnf, logger)
	defer sqlDB.Close()
	
	//database := db.New(sqlDB)

	app := application.NewApplication(ctx)
	r := app.NewRouter()
	
	server := http.Server{
		Addr:    ":" + cnf.Port,
		Handler: r,
	}
	return &server
}

func prepareDatabase(cnf *config.Config, logger log.Logger) *sql.DB {
	db, err := sql.Open("mysql", cnf.GetMySQLDSN())
	if err != nil {
		for range cnf.RetryCount {
			logger.Warn("failed to connect to database, retrying...", log.Err(err))
			time.Sleep(time.Second * 5)
			db, err = sql.Open("mysql", cnf.GetMySQLDSN())
			if err == nil {
				break
			}
		}
		if err != nil {
			logger.Panic("failed to connect to database", log.Err(err))
		}
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}