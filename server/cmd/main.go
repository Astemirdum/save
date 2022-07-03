package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Astemirdum/save/server/config"
	"github.com/Astemirdum/save/server/internal/handler"
	"github.com/Astemirdum/save/server/internal/repository"
	"github.com/Astemirdum/save/server/internal/service"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	rand.Seed(time.Now().UnixNano())

	if err := godotenv.Load("server.env"); err != nil {
		log.Fatalf("load envs from .env  %v", err)
	}

	// config
	cfg := config.NewConfig()

	// repo -> service -> handler
	repo := repository.NewRepository(cfg.FilePath)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services, log)

	addr1 := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port1)
	srv1 := handler.NewServer(handlers.NewRouter1(), addr1)
	go func() {
		if err := srv1.Run(); err != nil {
			log.Fatalf("Server init %v", err)
		}
	}()
	log.Printf("Server1 has been started on %s", addr1)

	addr2 := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port2)
	srv2 := handler.NewServer(handlers.NewRouter2(), addr2)
	go func() {
		if err := srv2.Run(); err != nil {
			log.Fatalf("Server init %v", err)
		}
	}()

	log.Printf("Server2 has been started on %s", addr2)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Graceful shutdown")
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFn()
	go func() {
		if err := srv1.Shutdown(ctx); err != nil {
			log.Errorf("Server Shutdown %v", err)
		}
	}()
	if err := srv1.Shutdown(ctx); err != nil {
		log.Errorf("Server Shutdown %v", err)
	}
}
