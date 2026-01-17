// accounts api
package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/andersonreyes/moneybadger/models"
	"github.com/andersonreyes/moneybadger/store"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type accountsApi struct {
	db store.AccountStore
}

func AccountsInit(ctx *context.Context, db *gorm.DB) accountsApi {
	return accountsApi{
		db: store.AccountStoreInit(ctx, db),
	}
}

type accountsRequestParams struct {
	AccountNumber string `uri:"accountNumber"`
}

func (c accountsApi) SetupRouter(router *gin.RouterGroup) error {
	router.GET("/accounts", func(ctx *gin.Context) {
		accs, err := c.db.ListAccounts()

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

	router.PUT("/accounts/:accountNumber", func(ctx *gin.Context) {

		var req accountsRequestParams

		if err := ctx.ShouldBindUri(&req); err != nil {
			log.Printf("invalid account number: %s\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		acc, err := c.db.GetAccount(req.AccountNumber)

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
		acc, err := c.db.GetAccount(req.AccountNumber)

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

		if err := ctx.ShouldBindJSON(&accToCreate); err != nil {
			log.Printf("Post:/accounts error: %s", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}

		if err := c.db.CreateAccount(accToCreate); err != nil {
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
		err := c.db.DeleteAccount(req.AccountNumber)

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

	return nil
}
