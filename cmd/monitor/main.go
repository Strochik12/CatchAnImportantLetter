package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Strochik12/CatchAnImportantLetter/internal/config"
	"github.com/Strochik12/CatchAnImportantLetter/internal/processor"
)

func main() {
	log.Println("📧 HSE Email Alert System")
	log.Println("==========================")

	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	log.Println("Конфигурация загружена")

	proc, err := processor.NewProcessor(cfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации системы: %v", err)
	}
	log.Println("Система инициализирована")

	// Настраиваем Gracefull Shotdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Таймер для переодического вывода статистики
	statsTicker := time.NewTicker(5 * time.Minute)
	defer statsTicker.Stop()

	go func() {
		for {
			select {
			case <-statsTicker.C:
				proc.PrintStats()
			case <-ctx.Done():
				return
			}
		}
	}()

	// Запускаем систему
	log.Println("🚀 Запускаем мониторинг почты...")
	log.Println("----------------------------------------")

	if err := proc.Start(ctx); err != nil {
		log.Fatalf("Ошибка работы системы: %v", err)
	}

	log.Println("\n----------------------------------------")
	log.Println("Система остановлена")
	proc.PrintStats()
}
