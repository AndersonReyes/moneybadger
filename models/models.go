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

type Account struct {
	gorm.Model
	AccountNumber  string          `gorm:"index;unique"`
	Type           AccountType     `gorm:"size:128;notnull"`
	Name           string          `gorm:"index;notnull"`
	InitialBalance decimal.Decimal `gorm:"notnull"`
	Balance        decimal.Decimal `gorm:"notnull"`
}

type ApiRoute interface {
	SetupRouter(router *gin.RouterGroup) error
}

type Transaction struct {
	gorm.Model
	Description          string          `gorm:"notnull;class:FULLTEXT"`
	Amount               decimal.Decimal `gorm:"notnull"`
	Category             string          `gorm:"index;notnull"`
	Date                 time.Time       `gorm:"index,notnull"`
	SourceAccountID      uint
	SourceAccount        Account
	DestinationAccountID uint
	DestinationAccount   Account
	Tags                 string
}

type TransactionFilters struct {
	TextSearch string    `uri:"textSearch"`
	StartDate  time.Time `uri:"startDate" binding:"required"`
	EndDate    time.Time `uri:"endDate" binding:"required"`
}

type Budget struct {
	gorm.Model
	Category       string `gorm:"index;notnull;unique"`
	ExpectedAmount decimal.Decimal
	ActualAmount   decimal.Decimal
}
