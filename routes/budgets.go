// budgets api
package routes

import (
	"log"
	"net/http"

	"github.com/andersonreyes/moneybadger/models"
	"github.com/andersonreyes/moneybadger/store"
	"github.com/gin-gonic/gin"
)

type budgetsApi struct {
	db store.Store
}

func BudgetsInit(dbStore store.Store) budgetsApi {
	return budgetsApi{
		db: dbStore,
	}
}

func (c budgetsApi) SetupRouter(router *gin.RouterGroup) error {

	router.GET("/budgets", func(ctx *gin.Context) {
		var filters models.TransactionFilters

		if err := ctx.ShouldBindQuery(&filters); err != nil {
			log.Printf("invalid query params: %s", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		budgets, err := c.db.Budgets.ListBudgets(filters)

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
			"data":    budgets,
		})
	})

	router.DELETE("/budgets/all", func(ctx *gin.Context) {

		err := c.db.Budgets.DeleteAllBudget()

		if err != nil {
			log.Printf("error deleting budgets: %s\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		ctx.Header("HX-Refresh", "true")
		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "all budgets deleted",
		})
	})

	router.POST("/budgets", func(ctx *gin.Context) {
		var budget models.Budget

		if err := ctx.ShouldBindJSON(&budget); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		if err := c.db.Budgets.CreateBudget(budget); err != nil {

		}

	})

	return nil
}
