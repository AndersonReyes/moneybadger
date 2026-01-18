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

type templateRequestParams struct {
	Name string
	Body map[string]any
}

func setupTemplatesApi(router *gin.RouterGroup) {

	router.POST("/template", func(ctx *gin.Context) {

		var body templateRequestParams

		if err := ctx.ShouldBindJSON(&body); err != nil {
			log.Printf("invalid template params: %s\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return

		}

		ctx.HTML(http.StatusOK, body.Name, body.Body)
	})
}

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

	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("views/templates/**/*")

	setupTemplatesApi(apiRouter)

	if err := routes.SetupRoutes(apiRouter, &ctx, db); err != nil {
		log.Panicf("failed to setup routes: %s\n", err)
	}

	if err := views.SetUpViews(router, store.AccountStoreInit(&ctx, db)); err != nil {
		log.Panicf("failed to setup views: %s\n", err)
	}

	router.Run()
}
