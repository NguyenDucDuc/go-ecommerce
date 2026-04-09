package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Transaction là một function nhận vào sessCtx và trả về error
type TransactionFunc func(sessCtx context.Context) error

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn TransactionFunc) error
}

type mongoTransactionManager struct {
	client *mongo.Client
}

func NewTransactionManager(client *mongo.Client) TransactionManager {
	return &mongoTransactionManager{client: client}
}

func (m *mongoTransactionManager) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	session, err := m.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		return nil, fn(sessCtx)
	})

	return err
}