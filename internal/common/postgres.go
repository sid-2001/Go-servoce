package common

import (
	"fmt"
	"os/exec"
	"strings"
)

// PostgresStore executes SQL commands through the local psql client.
// This keeps the example dependency-light while still using PostgreSQL as the only datastore.
type PostgresStore struct {
	DatabaseURL string
}

func NewPostgresStore(databaseURL string) *PostgresStore {
	return &PostgresStore{DatabaseURL: databaseURL}
}

// Exec runs a SQL statement that does not return rows.
func (p *PostgresStore) Exec(sql string) error {
	cmd := exec.Command("psql", p.DatabaseURL, "-v", "ON_ERROR_STOP=1", "-c", sql)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("psql exec failed: %w: %s", err, string(out))
	}
	return nil
}

// QueryJSON runs SQL and returns compact JSON output from PostgreSQL.
func (p *PostgresStore) QueryJSON(sql string) ([]byte, error) {
	cmd := exec.Command("psql", p.DatabaseURL, "-At", "-v", "ON_ERROR_STOP=1", "-c", sql)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("psql query failed: %w: %s", err, string(out))
	}
	return []byte(strings.TrimSpace(string(out))), nil
}

// QuoteLiteral safely wraps a string for SQL using standard single-quote escaping.
func QuoteLiteral(value string) string {
	escaped := strings.ReplaceAll(value, "'", "''")
	return "'" + escaped + "'"
}
