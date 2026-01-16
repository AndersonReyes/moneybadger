package models

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountType = string

const (
	AccountTypeDefault    = " default"
	AccountTypeCreditCard = "creditcard"
)

type Account struct {
	gorm.Model
	ID             uint            `gorm:"primaryKey;autoIncrement:true"`
	Type           AccountType     `gorm:"size:128"`
	AccountNumber  string          `gorm:"size:32;index;unique"`
	Name           string          `gorm:"index"`
	InitialBalance decimal.Decimal `gorm:"size:64"`
	Amount         decimal.Decimal `gorm:"size:64"`
}

type ApiRoute interface {
	SetupRouter(router *gin.RouterGroup) error
}
