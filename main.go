package main

import (
	"context"
	"log"

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
	db.AutoMigrate(&models.Account{})

	router := gin.Default()
	apiRouter := router.Group("/api")

	if err := routes.SetupRoutes(apiRouter, &ctx, db); err != nil {
		log.Panicf("failed to setup routes: %s\n", err)
	}

	if err := views.SetUpViews(router, store.AccountStoreInit(&ctx, db)); err != nil {
		log.Panicf("failed to setup views: %s\n", err)
	}
	router.Run()
}
