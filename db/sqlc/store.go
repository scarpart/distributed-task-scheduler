package db

import (
	"context"
	"database/sql"
	"fmt"
)

// `Store` provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// Creates and returns a new Store 
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

// Executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("tx err: %s, rb err: %s", err, rerr)
		}
		return err
	}

	err = tx.Commit()
	return err
}

// TODO: figure out what sorts of transactions I need to put in here

