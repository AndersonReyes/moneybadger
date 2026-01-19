package main

import (
	"context"
	"log"
	"net/http"

	"github.com/andersonreyes/moneybadger/models"
	"github.com/andersonreyes/moneybadger/routes"
	"github.com/andersonreyes/moneybadger/store"
	"github.com/andersonreyes/moneybadger/views"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	log.SetPrefix("money.server: ")
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database")
	}

	ctx := context.Background()

	// // Migrate the schema
	db.AutoMigrate(&models.Account{}, &models.Transaction{})

	router := gin.Default()
	apiRouter := router.Group("/api")

	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("views/templates/**/*")
	dbStore := store.StoreInit(&ctx, db)

	if err := routes.SetupRoutes(apiRouter, dbStore); err != nil {
		log.Panicf("failed to setup routes: %s\n", err)
	}

	if err := views.SetUpViews(router, dbStore); err != nil {
		log.Panicf("failed to setup views: %s\n", err)
	}

	apiRouter.GET("/template/raw", func(ctx *gin.Context) {
		var templateName = ctx.Query("templateName")

		if templateName == "" {
			log.Printf("invalid template name: %s\n", err)
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		ctx.HTML(http.StatusOK, templateName, gin.H{})
	})

	router.Run()
}
