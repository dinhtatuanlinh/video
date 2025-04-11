package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	//CreateOperatorTx(ctx context.Context, arg CreateOperatorTxParams) (CreateOperatorTxResult, error)
	//VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
	//ResendVerifyEmailTx(ctx context.Context, params ResendVerifyEmailTxParams) error
	//CreateAdminTx(ctx context.Context, arg CreateOperatorTxParams) (CreateOperatorTxResult, error)
}

type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
