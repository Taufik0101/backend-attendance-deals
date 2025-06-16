package repositories

import (
	"context"

	"gorm.io/gorm"
)

type ctxKey string

const (
	txKey ctxKey = "transaction-context"
)

func InjectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func ExtractTx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if ok {
		return tx
	}
	return defaultDB
}
