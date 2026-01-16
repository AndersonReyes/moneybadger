package views

import (
	"log"
	"net/http"

	"github.com/andersonreyes/moneybadger/store"
	"github.com/gin-gonic/gin"
)

func SetUpViews(router *gin.Engine, accountsStore store.AccountStore) error {
	router.LoadHTMLGlob("views/**/*.tmpl")
	router.GET("/accounts", func(ctx *gin.Context) {
		accs, err := accountsStore.ListAccounts()

		if err != nil {
			log.Printf("failed to get accounts: %s", err)
			ctx.JSON(500, gin.H{
				"error":   true,
				"message": err.Error(),
			})
		}

		ctx.HTML(http.StatusOK, "accounts.tmpl", gin.H{
			"title": accs,
		})
	})

	return nil
}
