package store

import (
	"context"
	"log"
	"time"

	"github.com/andersonreyes/moneybadger/models"
	"gorm.io/gorm"
)

type TransactionStore struct {
	ctx *context.Context
	db  *gorm.DB
}

func TransactionStoreInit(ctx *context.Context, db *gorm.DB) TransactionStore {
	return TransactionStore{
		ctx: ctx,
		db:  db,
	}
}

func (c TransactionStore) CreateTransaction(a models.Transaction) error {
	err := gorm.G[models.Transaction](c.db).Create(*c.ctx, &a)
	if err != nil {
		log.Println("CreateTransaction failed: " + err.Error())
	}
	return err
}

func (c TransactionStore) ListTransactions() ([]models.Transaction, error) {
	accounts, err := gorm.G[models.Transaction](c.db).Order("date desc").Find(*c.ctx)
	if err != nil {
		log.Println("ListTransactions failed: " + err.Error())
		return nil, err
	}
	return accounts, err
}

func (c TransactionStore) ListAccountTransactions(acc models.Account, startDate time.Time, endDate time.Time) ([]models.Transaction, error) {
	accounts, err := gorm.G[models.Transaction](c.db).
		Where("source_account_id = ? and date >=   and date < ?", acc.ID, startDate, endDate).
		Order("date desc").Find(*c.ctx)
	if err != nil {
		log.Println("ListTransactions failed: " + err.Error())
		return nil, err
	}
	return accounts, err
}

func (c TransactionStore) GetTransaction(id uint) (models.Transaction, error) {
	acc, err := gorm.G[models.Transaction](c.db).Where("id = ?", id).First(*c.ctx)
	if err != nil {
		log.Println("GetTransaction failed: " + err.Error())
		return models.Transaction{}, err
	}
	return acc, nil
}

func (c TransactionStore) UpdateTransaction(m models.Transaction) error {
	_, err := gorm.G[models.Transaction](c.db).Where("id = ?", m.ID).Updates(*c.ctx, m)

	if err != nil {
		log.Println("UpdateTransaction failed: " + err.Error())
		return err
	}
	return nil
}

func (c TransactionStore) DeleteTransaction(id uint) error {
	n, err := gorm.G[models.Transaction](c.db).Where("id = ?", id).Delete(*c.ctx)

	if n != 1 || err != nil {
		log.Printf("failed to delete account[%d]: %s\n", id, err)
		return err
	}

	return nil
}
