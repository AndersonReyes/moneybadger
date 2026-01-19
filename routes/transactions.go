// transactions api
package routes

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/andersonreyes/moneybadger/models"
	"github.com/andersonreyes/moneybadger/store"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type transactionsApi struct {
	db store.Store
}

func TransactionsInit(dbStore store.Store) transactionsApi {
	return transactionsApi{
		db: dbStore,
	}
}

type transactionsRequestParams struct {
	AccountNumber string `uri:"accountNumber"`
	TemplateName  string `uri:"templateName"`
}

func (c transactionsApi) importCsv(f *multipart.File) error {
	reader := csv.NewReader(*f)
	// skip header
	reader.Read()
	for {

		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if len(line) != 7 {
			log.Printf("invalid csv line, expected 7 fields but got [%d]: %+v\n", len(line), line)
			return errors.New("invalid csv line")
		}

		desc := line[4]

		var transaction models.Transaction

		amnt, err := decimal.NewFromString(line[3])
		if err != nil {
			log.Printf("invalid account column [%s]: %s\n", line[3], err.Error())
			return err
		}

		if line[1] != "" {
			sourceAccount, err := c.db.Accounts.GetAccountByName(line[1])
			if err == nil {
				transaction.SourceAccount = sourceAccount
				transaction.SourceAccountID = sourceAccount.ID
			}

		}

		if line[6] != "" {
			destAccount, err := c.db.Accounts.GetAccountByName(line[6])
			if err == nil {
				transaction.DestinationAccount = destAccount
				transaction.DestinationAccountID = destAccount.ID
			}

		}

		t, err := time.Parse(time.DateOnly, line[0])
		if err != nil {
			log.Printf("invalid date [%s] : %s\n", line[0], err.Error())
			return err
		}

		transaction.Date = t
		transaction.Amount = amnt
		transaction.Category = line[2]
		transaction.Description = desc

		_, found, err := c.db.Transactions.GetExistingTransaction(transaction)

		if err != nil {
			log.Printf("error looking for existing transaction [%+v]: %s", transaction, err.Error())
			return err

		}

		if !found {
			if err := c.db.Transactions.CreateTransaction(transaction); err != nil {
				log.Printf("failed to create transaction [%+v]: %s", transaction, err.Error())
				return err
			}
		}

	}

	return nil
}

func (c transactionsApi) SetupRouter(router *gin.RouterGroup) error {
	router.POST("/transactions/upload", func(ctx *gin.Context) {
		// Source
		file, err := ctx.FormFile("transactions")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
		}

		f, _ := file.Open()
		defer f.Close()

		if err := c.importCsv(&f); err != nil {
			ctx.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}

		ctx.Header("HX-Refresh", "true")
		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": fmt.Sprintf("File %s uploaded successfully.", file.Filename),
		})
	})

	router.GET("/transactions", func(ctx *gin.Context) {
		accs, err := c.db.Transactions.ListTransactions()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "",
			"data":    accs,
		})
	})

	router.DELETE("/transactions/all", func(ctx *gin.Context) {

		err := c.db.Transactions.DeleteAllTransaction()

		if err != nil {
			log.Printf("error deleting transactions: %s\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		ctx.Header("HX-Refresh", "true")
		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "all transactions deleted",
		})
	})

	return nil
}
