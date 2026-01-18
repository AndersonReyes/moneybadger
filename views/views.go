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

func SetUpViews(router *gin.Engine, accountsStore store.AccountStore) error {
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

	router.GET("/accounts/edit/:accountNumber", func(ctx *gin.Context) {
		var req accountViewRequestParams

		if err := ctx.ShouldBindUri(&req); err != nil || req.AccountNumber == "" {
			log.Printf("invalid account number [%s]: %s\n", req.AccountNumber, err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		acc, err := accountsStore.GetAccount(req.AccountNumber)
		if err != nil {
			log.Printf("failed to get account[%s] for editing: %s", req.AccountNumber, err)
			ctx.JSON(500, gin.H{
				"error":   true,
				"message": err.Error(),
			})
		}

		ctx.HTML(http.StatusOK, "accounts/edit.tmpl", acc)
	})

	return nil
}
