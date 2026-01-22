package store

import (
	"context"
	"log"

	"github.com/andersonreyes/moneybadger/models"
	"gorm.io/gorm"
)

type BudgetStore struct {
	ctx *context.Context
	db  *gorm.DB
}

func BudgetStoreInit(ctx *context.Context, db *gorm.DB) BudgetStore {
	return BudgetStore{
		ctx: ctx,
		db:  db,
	}
}

func (c BudgetStore) CreateBudget(a models.Budget) error {
	err := gorm.G[models.Budget](c.db).Create(*c.ctx, &a)
	if err != nil {
		log.Println("CreateBudget failed: " + err.Error())
	}
	return err
}

func (c BudgetStore) ListBudgets(filters models.TransactionFilters) ([]models.Budget, error) {
	q := `with t as (
	select category, amount from transactions where date >= ? and date < ?
	),
	actual as (
	select category, SUM(amount) as actual from t
	)

	select budgets.category,budgets.expected_amount, actual as actual_amount from budgets
	left join actual
		on budgets.category = actual.category

	order by budgets.category
	`
	budgets, err := gorm.G[models.Budget](c.db).Raw(q, filters.StartDate, filters.EndDate).
		Find(*c.ctx)

	if err != nil {
		log.Println("ListBudgets failed: " + err.Error())
		return nil, err
	}
	return budgets, err
}

func (c BudgetStore) GetBudget(category string) (models.Budget, error) {
	acc, err := gorm.G[models.Budget](c.db).
		Where("category = ?", category).
		First(*c.ctx)
	if err != nil {
		log.Println("GetBudget failed: " + err.Error())
		return models.Budget{}, err
	}
	return acc, nil
}

func (c BudgetStore) UpdateBudget(m models.Budget) error {
	_, err := gorm.G[models.Budget](c.db).Where("id = ?", m.ID).Updates(*c.ctx, m)

	if err != nil {
		log.Println("UpdateBudget failed: " + err.Error())
		return err
	}
	return nil
}

func (c BudgetStore) DeleteBudget(id uint) error {
	n, err := gorm.G[models.Budget](c.db).Where("id = ?", id).Delete(*c.ctx)

	if n != 1 || err != nil {
		log.Printf("failed to delete budget[%d]: %s\n", id, err)
		return err
	}

	return nil
}

func (c BudgetStore) DeleteAllBudget() error {
	n, err := gorm.G[models.Budget](c.db).Where(" 1 = 1", nil).Delete(*c.ctx)

	if n != 1 || err != nil {
		log.Printf("failed to delete all budgets %s\n", err)
		return err
	}

	return nil
}
