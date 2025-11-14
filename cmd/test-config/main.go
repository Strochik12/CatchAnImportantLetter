package main

import (
	"fmt"
	"log"

	"github.com/Strochik12/CatchAnImportantLetter/internal/config"
)

func main() {
	cfg, err := config.Load("configs/config.example.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Loaded configuration:\n")
	fmt.Printf("- IMAP Server: %s\n", cfg.IMAP.Server)
	fmt.Printf("- Rules count: %d\n", len(cfg.Rules))

	for _, rule := range cfg.Rules {
		fmt.Printf("  - %s (priority: %d)\n", rule.Name, rule.Priority)
	}
}
