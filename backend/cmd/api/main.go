package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nicolas/finanzas/backend/internal/config"
	"github.com/nicolas/finanzas/backend/internal/server"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	r := server.New()

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Server failed:", err)
	}
}