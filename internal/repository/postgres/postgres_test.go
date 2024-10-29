package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/juggler434/auth-server/pkg/models"
	"github.com/pashagolub/pgxmock/v4"
)

type AnyUUID struct{}

func (a AnyUUID) Match(v interface{}) bool {
	_, ok := v.(uuid.UUID)
	return ok
}

func TestCreateUser(t *testing.T) {
	tests := map[string]struct {
		setupMock          func(t *testing.T, mock pgxmock.PgxPoolIface)
		assertExpectations func(t *testing.T, mock pgxmock.PgxPoolIface, err error)
	}{
		"commits a new user to the database": {
			setupMock: func(t *testing.T, mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(AnyUUID{}, "testMan", "testMan@example.com", []byte("gobbledyGook")).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			assertExpectations: func(t *testing.T, mock pgxmock.PgxPoolIface, err error) {
				if err != nil {
					t.Fatalf("expected error to be nil, got %s", err)
				}
				if em := mock.ExpectationsWereMet(); em != nil {
					t.Fatalf("expected all sql expectation to be met, got %s", em)
				}
			},
		},
		"returns an error if there is an error commiting to the database": {
			setupMock: func(t *testing.T, mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(AnyUUID{}, "testMan", "testMan@example.com", []byte("gobbledyGook")).
					WillReturnError(fmt.Errorf("duplicate email"))
			},
			assertExpectations: func(t *testing.T, mock pgxmock.PgxPoolIface, err error) {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}

				if em := mock.ExpectationsWereMet(); em != nil {
					t.Fatalf("expected all sql expectation to be met, got %s", em)
				}
			},
		},
		"returns an error if number of rows affected is 0": {
			setupMock: func(t *testing.T, mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(AnyUUID{}, "testMan", "testMan@example.com", []byte("gobbledyGook")).
					WillReturnResult(pgxmock.NewResult("INSERT", 0))
			},
			assertExpectations: func(t *testing.T, mock pgxmock.PgxPoolIface, err error) {
				if em := mock.ExpectationsWereMet(); em != nil {
					t.Fatalf("expected all sql expectation to be met, got %s", em)
				}

				if err == nil {
					t.Fatal("expected an error, got nil")
				}
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("error getting mock DB connection %s", err)
			}

			test.setupMock(t, mock)
			db := postgresDb{
				dbPool: mock,
			}

			err = db.CreateUser(context.Background(), &models.User{
				ID:             uuid.New(),
				UserName:       "testMan",
				Email:          "testMan@example.com",
				HashedPassword: []byte("gobbledyGook"),
			})

			test.assertExpectations(t, mock, err)

		})
	}
}
