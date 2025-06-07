package database

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	config "github.com/vnnyx/employee-management/config/api"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

type PostgresConfig struct {
	Host                string
	Port                int64
	User                string
	Password            string
	DBName              string
	SSLMode             string
	MaxConn             int64
	EnableObservability bool
}

func GetPostgresConnection(config PostgresConfig) (*pgxpool.Pool, error) {
	if pool != nil {
		return pool, nil
	}

	connString := "host=" + config.Host +
		" port=" + strconv.FormatInt(config.Port, 10) +
		" user=" + config.User +
		" password=" + config.Password +
		" dbname=" + config.DBName +
		" sslmode=" + config.SSLMode

	if config.MaxConn > 0 {
		connString += " pool_max_conns=" + fmt.Sprint(config.MaxConn)
	}

	dbConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	if config.EnableObservability {
		dbConfig.ConnConfig.Tracer = otelpgx.NewTracer(
			otelpgx.WithTracerAttributes(semconv.DBNameKey.String(config.DBName)),
			otelpgx.WithTracerAttributes(semconv.DBSystemPostgreSQL),
		)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	pool = db

	return pool, nil
}

func GetPostgresConnectionSeed(config PostgresConfig) (*sqlx.DB, error) {

	connString := "host=" + config.Host +
		" port=" + strconv.FormatInt(config.Port, 10) +
		" user=" + config.User +
		" password=" + config.Password +
		" dbname=" + config.DBName +
		" sslmode=" + config.SSLMode

	db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetPool() *pgxpool.Pool {
	if pool == nil {
		once.Do(func() {
			config := config.Get()
			if config == nil {
				panic("Postgres configuration is not set")
			}

			db, err := GetPostgresConnection(PostgresConfig{
				Host:                config.Postgres.Host,
				Port:                config.Postgres.Port,
				User:                config.Postgres.User,
				Password:            config.Postgres.Password,
				DBName:              config.Postgres.DBName,
				SSLMode:             config.Postgres.SSLMode,
				MaxConn:             config.Postgres.MaxConn,
				EnableObservability: *config.Observability.Enable,
			})
			if err != nil {
				panic(errors.Wrap(err, "failed to initialize Postgres connection pool"))
			}

			pool = db
		})
	}
	return pool
}

func Atomic(ctx context.Context, txOpt pgx.TxOptions, fn func(tx DBTx) error) error {
	conn, err := GetPool().Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, txOpt)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("cannot rollback %w: %w", rbErr, err)
		}

		return errors.Wrap(err, "cb()")

	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "tx.Commit()")
	}
	return nil
}

func WithAuditContext(ctx context.Context, authCredential authCredential.Credential, txOpt pgx.TxOptions, fn func(tx DBTx) error) error {
	conn, err := GetPool().Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, txOpt)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, fmt.Sprintf(`SET LOCAL "app.current_user" = '%s';`, authCredential.UserID))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, fmt.Sprintf(`SET LOCAL "app.request_id" = '%s';`, authCredential.RequestID))
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("cannot rollback %w: %w", rbErr, err)
		}

		return errors.Wrap(err, "cb()")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "tx.Commit()")
	}

	return nil
}
