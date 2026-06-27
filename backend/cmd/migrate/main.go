package main

import (
	"fmt"
	"os"

	"github.com/nicolas/finanzas/backend/internal/config"
	"github.com/nicolas/finanzas/backend/internal/db"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate [up]")
		os.Exit(1)
	}

	cfg := config.Load()

	switch os.Args[1] {
	case "up":
		if err := db.Migrate(cfg); err != nil {
			fmt.Printf("Migration up failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations applied successfully")
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}