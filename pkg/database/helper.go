package database

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

func Rebind(query string) string {
	return sqlx.Rebind(sqlx.DOLLAR, query)
}

func In(query string, args ...any) (string, []any, error) {
	q, a, err := sqlx.In(query, args...)
	if err != nil {
		return "", nil, err
	}

	return Rebind(q), a, nil
}

func IsUniqueViolation(err error, constraintName string) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			if constraintName == "" || pgErr.ConstraintName == constraintName {
				return true
			}
		}
	}

	return false
}

func ExecuteSQLFile(db *sqlx.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading SQL file %s: %w", filePath, err)
	}

	query := strings.TrimSpace(string(content))
	if query == "" {
		return nil
	}

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("error executing query in file %s: %w", filePath, err)
	}

	return nil
}
