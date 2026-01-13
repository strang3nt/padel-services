package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const databaseTablesPath = "resources/sql/create_tables.sql"

func CreateDatabaseTables(ctx context.Context, conn *pgxpool.Pool) error {

	databaseTables, err := os.ReadFile(databaseTablesPath)
	if err != nil {
		return fmt.Errorf("unable to read database tables file: %w", err)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("error while rolling back transaction: %v", err)
		}
	}()

	_, err = tx.Exec(ctx, string(databaseTables))
	if err != nil {
		return fmt.Errorf("unable to execute query: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("unable to create database table: %w", err)
	}
	return nil
}
