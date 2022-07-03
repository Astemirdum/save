package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/Astemirdum/save/client/config"
	"github.com/Astemirdum/save/client/service"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	log := zap.NewExample().Named("client")
	if err := godotenv.Load("client.env"); err != nil {
		log.Fatal("load envs from .env  ", zap.Error(err))
	}

	// config
	cfg := config.NewConfig()

	svc := service.NewFileClientService(cfg, log.Named("service"))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := svc.PostSave(ctx); err != nil {
			log.Fatal("save req fail", zap.Error(err))
		}
	}()
	go func() {
		if err := svc.PutWrite(ctx); err != nil {
			log.Fatal("save req fail", zap.Error(err))
		}
	}()
	go func() {
		if err := svc.GetFileCount(ctx); err != nil {
			log.Fatal("save req fail", zap.Error(err))
		}
	}()
	go func() {
		if err := svc.GetSrvTime(ctx); err != nil {
			log.Fatal("save req fail", zap.Error(err))
		}
	}()
	go func() {
		if err := svc.GetFileText(ctx); err != nil {
			log.Fatal("save req fail", zap.Error(err))
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("Graceful shutdown")
	cancel()
}
