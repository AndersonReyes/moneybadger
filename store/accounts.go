package store

import (
	"context"
	"log"

	"github.com/andersonreyes/moneybadger/models"
	"gorm.io/gorm"
)

type AccountStore struct {
	ctx *context.Context
	db  *gorm.DB
}

func AccountStoreInit(ctx *context.Context, db *gorm.DB) AccountStore {
	return AccountStore{
		ctx: ctx,
		db:  db,
	}
}

func (c AccountStore) CreateAccount(a models.Account) error {
	err := gorm.G[models.Account](c.db).Create(*c.ctx, &a)
	if err != nil {
		log.Println("CreateAccount failed: " + err.Error())
	}
	return err
}

func (c AccountStore) ListAccounts() ([]models.Account, error) {
	accounts, err := gorm.G[models.Account](c.db).Order("name asc").Find(*c.ctx)
	if err != nil {
		log.Println("ListAccounts failed: " + err.Error())
		return nil, err
	}
	return accounts, err
}

func (c AccountStore) GetAccount(accountNumber string) (models.Account, error) {
	acc, err := gorm.G[models.Account](c.db).Where("account_number = ?", accountNumber).First(*c.ctx)
	if err != nil {
		log.Println("GetAccount failed: " + err.Error())
		return models.Account{}, err
	}
	return acc, nil
}
