package main

import (
	"context"
	"expvar"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/shaneu/indahaus/cmd/api/handlers"
	"github.com/shaneu/indahaus/pkg/auth"
	"github.com/shaneu/indahaus/pkg/database"

	"github.com/spf13/viper"
)

var build = "develop"

func main() {
	// for our usecase human readable logs instead of a structured logging solution will be adequate until and unless
	// we decide we want structured logging
	log := log.New(os.Stdout, "API: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

// run handles intitializing our app and will return an error in the case of failure
func run(log *log.Logger) error {
	// ===========================================================
	// Initialize configuration
	var cfg struct {
		Address   string
		Port      string
		DebugPort string
		DB        struct {
			Uri string
		}
		App struct {
			ReadTimeout     time.Duration
			ShutdownTimeout time.Duration
			WriteTimeout    time.Duration
		}
		Auth struct {
			Password string
			Username string
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	// supports overriding nested config fields with env vars, see
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "reading config")
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return errors.Wrap(err, "unmarshal config")
	}

	// Register `build` var with expvar so /debug/vars will reflect current build
	expvar.NewString("build").Set(build)

	log.Printf("main: Application initializing: version %q", build)
	defer log.Printf("main: completed")

	// ===========================================================
	// Initialize database
	log.Println("main: Initializing database support")
	db, err := database.Open(database.Config{
		Uri: cfg.DB.Uri,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer func() {
		log.Printf("main: Database Stopping")
		db.Close()
	}()

	// ===========================================================
	// Initialize debug endpoint
	// Not critical for application function so we do not abort startup or shutdown app if endpoints fails
	go func() {
		log.Printf("main: Debug Listening  :%s/debug/vars", cfg.DebugPort)

		if err := http.ListenAndServe(net.JoinHostPort(cfg.Address, cfg.DebugPort), http.DefaultServeMux); err != nil {
			log.Printf("main: debug listener closed: %v", err)
		}
	}()

	// channel to listen for SIGINT and SIGTERM signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// ===========================================================
	// Initialize auth
	a := auth.New(cfg.Auth.Username, cfg.Auth.Password)

	api := http.Server{
		Addr:         net.JoinHostPort(cfg.Address, cfg.Port),
		Handler:      handlers.API(build, a, db, log),
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("main: Api listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// select blocks until it either receives an error from listenAndServe or it receives a shut signal signal
	// it which cases it attempts to do a graceful shutdown
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Printf("main: %v : Start shutdown", sig)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.App.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
