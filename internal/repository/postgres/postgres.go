package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juggler434/auth-server/pkg/models"
)

// Used to wrap database functions to make testing easier

type postgresPool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type postgresDb struct {
	dbPool postgresPool
}

func newPostgresDb(ctx context.Context) (*postgresDb, error) {
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the postgres db: %w", err)
	}

	return &postgresDb{
		dbPool: pool,
	}, nil
}

func (db *postgresDb) CreateUser(ctx context.Context, user *models.User) error {
	res, err := db.dbPool.Exec(ctx, `
    INSERT INTO users (id, username, email, hashed_password) 
    VALUES ($1, $2, $3, $4)
    `, user.ID, user.UserName, user.Email, user.HashedPassword)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("no user inserted to database")
	}
	return nil
}
