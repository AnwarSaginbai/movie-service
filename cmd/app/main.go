package main

import (
	"context"
	"expvar"
	"github.com/AnwarSaginbai/netflix-service/internal/config"
	"github.com/AnwarSaginbai/netflix-service/internal/data"
	"github.com/AnwarSaginbai/netflix-service/internal/mailer"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"log/slog"
	"runtime"
	"sync"
	"time"
)

const version = "1.0.0"

type application struct {
	logger *slog.Logger
	config *config.Config
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	cfg := config.SetupConfig()
	logger := setupLogger(cfg.Env)
	moviesCollection, err := setupDB("movies")
	if err != nil {
		logger.Error("err", slog.String("db movies", err.Error()))
	}
	tokenCollection, err := setupDB("tokens")
	if err != nil {
		logger.Error("err", slog.String("db token", err.Error()))
	}
	usersCollection, err := setupDB("users")
	if err != nil {
		logger.Error("err", slog.String("db users", err.Error()))
	}
	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	var databaseStats bson.Raw
	err = moviesCollection.Database().RunCommand(
		context.Background(),
		bson.D{{Key: "dbStats", Value: 1}},
	).Decode(&databaseStats)
	if err != nil {
		logger.Error("err", slog.String("db stats", err.Error()))
	} else {
		expvar.Publish("database", expvar.Func(func() interface{} {
			return databaseStats.String()
		}))
	}

	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	logger.Info("database connected", slog.Bool("state", true))
	app := &application{
		logger: logger,
		config: cfg,
		models: data.NewModels(moviesCollection, usersCollection, tokenCollection),
		mailer: mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender),
	}
	if err := app.serve(); err != nil {
		log.Fatalln(err)
	}
}
