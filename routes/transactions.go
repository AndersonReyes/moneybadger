// transactions api
package routes

import (
	"context"

	"github.com/andersonreyes/moneybadger/store"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type transactionsApi struct {
	db store.TransactionStore
}

func TransactionsInit(ctx *context.Context, db *gorm.DB) transactionsApi {
	return transactionsApi{
		db: store.TransactionStoreInit(ctx, db),
	}
}

type transactionsRequestParams struct {
	AccountNumber string `uri:"accountNumber"`
	TemplateName  string `uri:"templateName"`
}

func (c transactionsApi) SetupRouter(router *gin.RouterGroup) error {
	router.GET("/transactions", func(ctx *gin.Context) {
		accs, err := c.db.ListTransactions()

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

	// router.PUT("/transactions", func(ctx *gin.Context) {

	// 	var overWrites models.Account

	// 	if err := ctx.ShouldBind(&overWrites); err != nil {
	// 		log.Printf("Put:/transactions error: %s", err)
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
	// 		return
	// 	}

	// 	if overWrites.Type == "" {
	// 		overWrites.Type = models.AccountTypeDefault
	// 	}

	// 	log.Printf("updating %+v\n", overWrites)
	// 	if err := c.db.UpdateAccount(overWrites); err != nil {
	// 		log.Printf("error updating account %+v:\n%s\n", overWrites, err)
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{
	// 			"error":   true,
	// 			"message": err.Error(),
	// 		})
	// 		return
	// 	}

	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"error":   false,
	// 		"message": "",
	// 	})
	// })

	// router.GET("/transactions/:accountNumber", func(ctx *gin.Context) {
	// 	var req transactionsRequestParams

	// 	if err := ctx.ShouldBindUri(&req); err != nil {
	// 		log.Printf("invalid account number: %s\n", err)
	// 		ctx.JSON(http.StatusBadRequest, gin.H{
	// 			"error":   true,
	// 			"message": err.Error(),
	// 		})
	// 		return
	// 	}
	// 	acc, err := c.db.GetAccount(req.AccountNumber)

	// 	if err != nil {
	// 		log.Printf("error getting account %+v:\n%s\n", req, err)
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{
	// 			"error":   true,
	// 			"message": err.Error(),
	// 		})
	// 		return
	// 	}

	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"error":   false,
	// 		"message": "",
	// 		"data":    []models.Account{acc},
	// 	})
	// })

	// router.POST("/transactions", func(ctx *gin.Context) {
	// 	var accToCreate models.Account
	// 	if accToCreate.Type == "" {
	// 		accToCreate.Type = models.AccountTypeDefault
	// 	}

	// 	if err := ctx.ShouldBind(&accToCreate); err != nil {
	// 		log.Printf("Post:/transactions error: %s", err)
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
	// 		return
	// 	}

	// 	if err := c.db.CreateAccount(accToCreate); err != nil {
	// 		log.Printf("error creating account %+v:\n%s\n", accToCreate, err)
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{
	// 			"error":   true,
	// 			"message": err.Error(),
	// 		})
	// 		return
	// 	}

	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"error":   false,
	// 		"message": "account created",
	// 	})
	// })

	// router.DELETE("/transactions/:accountNumber", func(ctx *gin.Context) {

	// 	var req transactionsRequestParams

	// 	if err := ctx.ShouldBindUri(&req); err != nil {
	// 		log.Printf("invalid account number: %s\n", err)
	// 		ctx.JSON(http.StatusBadRequest, gin.H{
	// 			"error":   true,
	// 			"message": err.Error(),
	// 		})
	// 		return
	// 	}
	// 	err := c.db.DeleteAccount(req.AccountNumber)

	// 	if err != nil {
	// 		log.Printf("error deleting account %s:\n%s\n", req.AccountNumber, err)
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{
	// 			"error":   true,
	// 			"message": err.Error(),
	// 		})
	// 		return
	// 	}

	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"error":   false,
	// 		"message": "account deleted",
	// 	})
	// })

	// // templates

	// router.GET("/template/transactions/:templateName/:accountNumber", func(ctx *gin.Context) {
	// 	var req transactionsRequestParams

	// 	if err := ctx.ShouldBindUri(&req); err != nil {
	// 		log.Printf("invalid template name: %s\n", err)
	// 		ctx.JSON(http.StatusBadRequest, gin.H{
	// 			"error":   true,
	// 			"message": err.Error(),
	// 		})
	// 		return
	// 	}

	// 	a, err := c.db.GetAccount(req.AccountNumber)

	// 	if err != nil {
	// 		ctx.JSON(http.StatusNotFound, gin.H{
	// 			"error":   true,
	// 			"message": "invalid account number",
	// 		})
	// 	}

	// 	ctx.HTML(http.StatusOK, req.TemplateName, a)
	// })

	return nil
}
