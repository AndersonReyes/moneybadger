// accounts api
package routes

import (
	"log"
	"net/http"

	"github.com/andersonreyes/moneybadger/models"
	"github.com/andersonreyes/moneybadger/store"
	"github.com/gin-gonic/gin"
)

type accountsApi struct {
	db store.Store
}

func AccountsInit(dbStore store.Store) accountsApi {
	return accountsApi{
		db: dbStore,
	}
}

type accountsRequestParams struct {
	AccountNumber string `uri:"accountNumber"`
	TemplateName  string `uri:"templateName"`
}

func (c accountsApi) SetupRouter(router *gin.RouterGroup) error {
	router.GET("/accounts", func(ctx *gin.Context) {
		accs, err := c.db.Accounts.ListAccounts()

		if err != nil {
			ctx.JSON(500, gin.H{
				"error":   true,
				"message": err.Error(),
			})
		}

		ctx.JSON(200, gin.H{
			"error":   false,
			"message": "",
			"data":    accs,
		})
	})

	router.PUT("/accounts", func(ctx *gin.Context) {

		var overWrites models.Account

		if err := ctx.ShouldBind(&overWrites); err != nil {
			log.Printf("Put:/accounts error: %s", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}

		if overWrites.Type == "" {
			overWrites.Type = models.AccountTypeDefault
		}

		log.Printf("updating %+v\n", overWrites)
		if err := c.db.Accounts.UpdateAccount(overWrites); err != nil {
			log.Printf("error updating account %+v:\n%s\n", overWrites, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "",
		})
	})

	router.GET("/accounts/:accountNumber", func(ctx *gin.Context) {
		var req accountsRequestParams

		if err := ctx.ShouldBindUri(&req); err != nil {
			log.Printf("invalid account number: %s\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		acc, err := c.db.Accounts.GetAccount(req.AccountNumber)

		if err != nil {
			log.Printf("error getting account %+v:\n%s\n", req, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "",
			"data":    []models.Account{acc},
		})
	})

	router.POST("/accounts", func(ctx *gin.Context) {
		var accToCreate models.Account
		if accToCreate.Type == "" {
			accToCreate.Type = models.AccountTypeDefault
		}

		if err := ctx.ShouldBind(&accToCreate); err != nil {
			log.Printf("Post:/accounts error: %s", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}

		if err := c.db.Accounts.CreateAccount(accToCreate); err != nil {
			log.Printf("error creating account %+v:\n%s\n", accToCreate, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "account created",
		})
	})

	router.DELETE("/accounts/:accountNumber", func(ctx *gin.Context) {

		var req accountsRequestParams

		if err := ctx.ShouldBindUri(&req); err != nil {
			log.Printf("invalid account number: %s\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		err := c.db.Accounts.DeleteAccount(req.AccountNumber)

		if err != nil {
			log.Printf("error deleting account %s:\n%s\n", req.AccountNumber, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "account deleted",
		})
	})

	// templates

	router.GET("/template/accounts/:templateName/:accountNumber", func(ctx *gin.Context) {
		var req accountsRequestParams

		if err := ctx.ShouldBindUri(&req); err != nil {
			log.Printf("invalid template name: %s\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		a, err := c.db.Accounts.GetAccount(req.AccountNumber)

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "invalid account number",
			})
		}

		ctx.HTML(http.StatusOK, req.TemplateName, a)
	})

	return nil
}
