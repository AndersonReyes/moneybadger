package models

import (
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
	AccountNumber  string          `gorm:"primaryKey;unique"`
	Type           AccountType     `gorm:"size:128;notnull"`
	Name           string          `gorm:"index;notnull"`
	InitialBalance decimal.Decimal `gorm:"size:64;notnull"`
	Balance        decimal.Decimal `gorm:"size:64;notnull"`
}

type ApiRoute interface {
	SetupRouter(router *gin.RouterGroup) error
}
