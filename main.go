package main

import (
	"fmt"
	"log"

	"github.com/adarsh-jaiss/go-bank/api"
	"github.com/adarsh-jaiss/go-bank/db"
	"github.com/adarsh-jaiss/go-bank/middleware"
	"github.com/adarsh-jaiss/go-bank/store"
	"github.com/gin-gonic/gin"
)

func main() {
	err := db.Connect()
	if err != nil {
		fmt.Errorf("Error connecting to database %s", err)
	}

	fmt.Println("Database connected successfully")

	defer db.Disconnect()

		if err := db.CreateTable(); err != nil {
			log.Fatal("Error creating table schema", err)
		}

		userStore := store.NewPostgresUserStore(db.DB)
	    userHandler := api.NewUserHandler(userStore)
		authHandler := api.NewAuthHandler(userStore)

		app := gin.Default()
		app.Group("api")
		appV1 := app.Group("api/v1")
		auth := app.Group("api/v1/login", middleware.JWTAuthentication(userStore))

	// Versioned API routes
	// This is user handlers
	appV1.POST("/signup/user", userHandler.HandlePostUser())
	auth.POST("/user", authHandler.HandleAuthentication)		//gonna use it as AUTHENTICATION

	log.Fatal(app.Run(":8080"))
}
