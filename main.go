package main

import (
	"internal-transfers/config"
	"internal-transfers/db"
	"internal-transfers/routes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	db.Init()

	router := gin.Default()
	routes.RegisterRoutes(router, db.Conn)

	server := &http.Server{
		Addr:           ":" + config.App.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
