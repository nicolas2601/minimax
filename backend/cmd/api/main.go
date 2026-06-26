package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/nicolas/finanzas/backend/internal/config"
	"github.com/nicolas/finanzas/backend/internal/db"
	"github.com/nicolas/finanzas/backend/internal/server"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	gormDB, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := server.New(gormDB)

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}