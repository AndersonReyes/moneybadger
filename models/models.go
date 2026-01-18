package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountType = string

const (
	AccountTypeDefault    = "default"
	AccountTypeCreditCard = "creditcard"
)

type Category struct {
	Name               string `gorm:"notnull"`
	TransactionMatcher string
}

type Group struct {
	Name       string
	Categories []Category
}

type Account struct {
	gorm.Model
	AccountNumber  string          `gorm:"primaryKey;unique" binding:"required"`
	Type           AccountType     `gorm:"size:128;notnull"`
	Name           string          `gorm:"index;notnull"`
	InitialBalance decimal.Decimal `gorm:"notnull"`
	Balance        decimal.Decimal `gorm:"notnull"`
}

type Transaction struct {
	gorm.Model
	Description          string          `gorm:"notnull;class:FULLTEXT"`
	Amount               decimal.Decimal `gorm:"notnull"`
	Category             Category        `gorm:"notnull"`
	Date                 time.Time       `gorm:"index,notnull"`
	SourceAccount        Account
	SourceAccountID      string
	DestinationAccount   Account
	DestinationAccountID string
	Tags                 []string
}

type ApiRoute interface {
	SetupRouter(router *gin.RouterGroup) error
}
