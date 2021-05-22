package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/shaneu/indahaus/internal/data/schema"
	"github.com/shaneu/indahaus/pkg/database"
	"github.com/spf13/viper"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "usage: admin [command]")
		os.Exit(2)
	}

	log := log.New(os.Stdout, "API: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {
	// ===========================================================
	// Initialize configuration
	var cfg struct {
		DB struct {
			Uri      string
			Username string
			Password string
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
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

	switch os.Args[1] {
	case "migrate":
		dbCfg := database.Config{
			Password: cfg.DB.Password,
			Uri:      cfg.DB.Uri,
			Username: cfg.DB.Username,
		}
		err := migrate(dbCfg)
		if err != nil {
			return err
		}
	default:
		log.Fatal("unsupported command")
	}

	return nil
}

func migrate(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "unable to open database")
	}

	if err := schema.Migrate(db); err != nil {
		return errors.Wrap(err, "unable to migrate database")
	}

	return nil
}
