package views

import (
	"log"
	"net/http"

	"github.com/andersonreyes/moneybadger/store"
	"github.com/gin-gonic/gin"
)

type accountViewRequestParams struct {
	AccountNumber string `uri:"accountNumber"`
}

func SetUpViews(router *gin.Engine, accountsStore store.AccountStore,
	transactionStore store.TransactionStore) error {
	router.GET("/accounts", func(ctx *gin.Context) {
		accs, err := accountsStore.ListAccounts()

		if err != nil {
			log.Printf("failed to get accounts: %s", err)
			ctx.JSON(500, gin.H{
				"error":   true,
				"message": err.Error(),
			})
		}

		ctx.HTML(http.StatusOK, "accounts/index.tmpl", gin.H{
			"accounts": accs,
		})
	})

	router.GET("/transactions", func(ctx *gin.Context) {
		ts, err := transactionStore.ListTransactions()

		if err != nil {
			log.Printf("failed to get accounts: %s", err)
			ctx.JSON(500, gin.H{
				"error":   true,
				"message": err.Error(),
			})
		}

		ctx.HTML(http.StatusOK, "transactions/index.tmpl", gin.H{
			"transactions": ts,
		})
	})

	return nil
}
