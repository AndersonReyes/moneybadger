package routes

import (
	"log"

	"github.com/andersonreyes/moneybadger/models"
	"github.com/andersonreyes/moneybadger/store"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(g *gin.RouterGroup, dbStore store.Store) error {

	apis := []models.ApiRoute{AccountsInit(dbStore), TransactionsInit(dbStore)}

	for _, api := range apis {
		if err := api.SetupRouter(g); err != nil {
			log.Printf("failed to setup Accounts api: %s\n", err)
			return err
		}
	}

	return nil
}
