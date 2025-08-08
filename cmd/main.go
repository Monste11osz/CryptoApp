package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testYTask/internal/app"
	"testYTask/internal/config"
	"testYTask/internal/domain/models"

	"go.uber.org/zap"
)

func main() {
	cfg := config.GetConfig()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	application := app.NewApp(cfg)
	if err := application.Init(ctx); err != nil {
		panic(fmt.Sprintf("Failed to start application %v", zap.Error(err)))
	}
	defer func() {
		if r := recover(); r != nil {
			report := models.GetPanicReport(2, r)
			log.Printf("[MAIN PANIC] %v", report)
			os.Exit(1)
		}
	}()
	application.Run(ctx)

}
