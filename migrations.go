package framework

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"path/filepath"
)

type Migration struct {
	db        *sql.DB
	tableName string
}

func NewMigration(db *sql.DB) *Migration {
	return &Migration{db: db, tableName: "_migrations"}
}
func (o *Migration) createMigrationsTableIfNotExists() error {
	_, err := o.db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (name TEXT PRIMARY KEY)`, o.tableName))
	return err
}

func (o *Migration) Migrate(migrationsDir string) error {
	if err := o.createMigrationsTableIfNotExists(); err != nil {
		return err
	}

	// Get list of applied migrations
	appliedMigrations, err := o.getAppliedMigrations()

	if err != nil {
		return err
	}

	// Get list of all .sql files in the migrations directory
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Skip if this migration has already been applied
		if _, ok := appliedMigrations[file.Name()]; ok {
			continue
		}

		// Read the migration file
		content, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
		if err != nil {
			return err
		}

		// Start a transaction
		tx, err := o.db.Begin()
		if err != nil {
			return err
		}

		// Execute the migration
		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return err
		}

		// Record that the migration has been applied
		if _, err := tx.Exec(fmt.Sprintf("INSERT INTO %s (name) VALUES ($1)", o.tableName), file.Name()); err != nil {
			tx.Rollback()
			return err
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (o *Migration) getAppliedMigrations() (map[string]struct{}, error) {
	rows, err := o.db.Query(fmt.Sprintf("SELECT name FROM %s", o.tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appliedMigrations := make(map[string]struct{})
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		appliedMigrations[name] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return appliedMigrations, nil
}
