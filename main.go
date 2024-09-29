package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/dubass83/go-micro-auth/cmd/api"
	data "github.com/dubass83/go-micro-auth/data/sqlc"
	"github.com/dubass83/go-micro-auth/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// var interaptSignals = []os.Signal{
// 	os.Interrupt,
// 	syscall.SIGTERM,
// 	syscall.SIGINT,
// }

func main() {
	conf, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot load configuration")
	}
	if conf.Enviroment == "devel" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		// log.Debug().Msgf("config values: %+v", conf)
	}

	// ctx, stop := signal.NotifyContext(context.Background(), interaptSignals...)
	// defer stop()

	connPool, err := pgxpool.NewWithConfig(context.Background(), PoolConfig(conf))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot validate db connection string")
	}
	store := data.NewStore(connPool)
	runDbMigration(conf.MigrationURL, conf.DBSource)
	runChiServer(conf, store)
}

// PoolConfig create config for db connection pool
func PoolConfig(conf util.Config) *pgxpool.Config {

	dbConfig, err := pgxpool.ParseConfig(conf.DBSource)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to create a config")
	}
	if conf.DBPoolMaxConns != 0 {
		dbConfig.MaxConns = conf.DBPoolMaxConns
	}
	if conf.DBPoolMinConns != 0 {
		dbConfig.MinConns = conf.DBPoolMinConns
	}
	if conf.DBPoolMaxConnLifetime != time.Second*0 {
		dbConfig.MaxConnLifetime = conf.DBPoolMaxConnLifetime
	}
	if conf.DBPoolMaxConnIdleTime != time.Second*0 {
		dbConfig.MaxConnIdleTime = conf.DBPoolMaxConnIdleTime
	}
	if conf.DBPoolHealthCheckPeriod != time.Second*0 {
		dbConfig.HealthCheckPeriod = conf.DBPoolHealthCheckPeriod
	}
	if conf.DBPoolConnectTimeout != time.Second*0 {
		dbConfig.ConnConfig.ConnectTimeout = conf.DBPoolConnectTimeout
	}

	return dbConfig
}

// runDbMigration run db migration from file to db
func runDbMigration(migrationURL, dbSource string) {
	m, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("method", "main").
			Msg("can not create migration instance")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().
			Err(err).
			Msg("can not run migration up")
	}
	log.Info().Msg("successfully run db migration")
}

// runChiServer run http server with Chi framework
func runChiServer(conf util.Config, store data.Store) {
	server := api.CreateNewServer(conf, store)

	server.ConfigureCORS()
	server.AddMiddleware()
	server.MountHandlers()
	log.Info().
		Msgf("start listening on the port %s\n", server.Config.HTTPAddressString)
	err := http.ListenAndServe(server.Config.HTTPAddressString, server.Router)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("method", "main").
			Msg("can not start server")
	}
}
