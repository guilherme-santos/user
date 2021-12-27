package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/guilherme-santos/user"
	"github.com/guilherme-santos/user/http"
	"github.com/guilherme-santos/user/mysql"
	"github.com/guilherme-santos/user/stub"

	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var cfg struct {
	Logger struct {
		Level string `envconfig:"LEVEL" default:"info"`
	} `envconfig:"LOGGER"`
	HTTP struct {
		Addr string `envconfig:"ADDR" default:"0.0.0.0:80"`
	} `envconfig:"HTTP"`
	MySQL struct {
		Host     string `envconfig:"HOST" required:"true"`
		Port     int    `envconfig:"PORT" default:"3306"`
		User     string `envconfig:"USER" required:"true"`
		Password string `envconfig:"PASSWORD" required:"true"`
		Database string `envconfig:"DATABASE" default:"user"`
	} `envconfig:"MYSQL"`
}

func init() {
	envconfig.MustProcess("usersvc", &cfg)
}

func main() {
	log := logrus.New()
	logLevel, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		logLevel = logrus.InfoLevel
		log.WithError(err).Error("unable to set log level")
	} else {
		log.SetLevel(logLevel)
		log.WithField("log_level", logLevel).Info("log level updated")
	}

	db, err := mysql.NewConnection(
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Database,
	)
	if err != nil {
		log.WithError(err).Fatal("unable to connect to database")
	}
	log.Info("connected to database")

	userstorage := mysql.NewUserStorage(db)
	eventsvc := stub.NewEventService()
	usersvc := user.NewService(userstorage, eventsvc)

	httprouter := http.NewRouter(log)
	http.NewHealthHandler(httprouter, db)
	httprouter.Route("/v1", func(r chi.Router) {
		http.NewUserHandler(r, usersvc)
	})

	httpsrv := http.NewServer(cfg.HTTP.Addr, httprouter)
	log.WithField("addr", cfg.HTTP.Addr).Info("running http server")

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		errCh <- httpsrv.ListenAndServe()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case err := <-errCh:
		log.WithError(err).Error("unable to run http server")
	case sig := <-sigCh:
		log.WithField("signal", sig).Warn("signal received, shuting down")
		// gracefullt shutdown http server
		httpsrv.Shutdown(context.Background())
	}
}
