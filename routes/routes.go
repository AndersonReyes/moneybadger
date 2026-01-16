package routes

import (
	"context"
	"log"

	"github.com/andersonreyes/moneybadger/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(g *gin.RouterGroup, ctx *context.Context, db *gorm.DB) error {

	apis := []models.ApiRoute{AccountsInit(ctx, db)}

	for _, api := range apis {
		if err := api.SetupRouter(g); err != nil {
			log.Printf("failed to setup Accounts api: %s\n", err)
			return err
		}
	}

	return nil
}
