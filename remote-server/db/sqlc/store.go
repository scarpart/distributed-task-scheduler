package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/scarpart/distributed-task-scheduler/util/logger"
)

// `Store` provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// Creates and returns a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return store.db.BeginTx(ctx, opts)
}

// Executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	logger.InfoLogger.Println("Executing transaction")

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		logger.ErrLogger.Printf("Error beginning transaction: %s\n", err.Error())
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		logger.ErrLogger.Printf("Error executing transaction callback: %s\n", err.Error())
		if rerr := tx.Rollback(); rerr != nil {
			logger.ErrLogger.Printf("Error on transaction rollback: %s\n", rerr.Error())
			return fmt.Errorf("tx err: %s, rb err: %s", err, rerr)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		logger.ErrLogger.Printf("Error committing transaction: %s\n", err.Error())
	} else {
		logger.InfoLogger.Println("Transaction successfully committed")
	}
	return err
}

// TODO: figure out what sorts of transactions I need to put in here
