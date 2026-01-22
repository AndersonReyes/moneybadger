package store

import (
	"context"

	"gorm.io/gorm"
)

type Store struct {
	Accounts     AccountStore
	Transactions TransactionStore
	Budgets      BudgetStore
}

func StoreInit(ctx *context.Context, db *gorm.DB) Store {
	return Store{
		Accounts:     AccountStoreInit(ctx, db),
		Transactions: TransactionStoreInit(ctx, db),
		Budgets:      BudgetStoreInit(ctx, db),
	}
}
