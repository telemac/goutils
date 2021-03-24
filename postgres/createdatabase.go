package postgres

import "fmt"

// CreateDatabaseIfNotExists returns the postgresql syntax equivalent to mysql CREATE DATABASE IF NOT EXISTS.
func CreateDatabaseIfNotExists(databaseName string) string {
	return fmt.Sprintf("SELECT 'CREATE DATABASE %s' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')\\gexec\n", databaseName, databaseName)
}
