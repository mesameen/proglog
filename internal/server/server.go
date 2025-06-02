package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mesameen/proglog/internal/config"
)

func StartHTTPServer(ctx context.Context) error {
	handler := newHandler()
	r := gin.Default()
	r.POST("/", handler.handleProduce)
	r.GET("/", handler.handleConsume)
	srv := http.Server{
		Addr:    config.CommonConfig.Port,
		Handler: r,
	}
	// starting server in seperate go rotine
	go func() {
		log.Printf("server is up and running on %v\n", config.CommonConfig.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicf("Failed to start server. Error: %v\n", err)
		}
	}()
	// for graceful shutdonw
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	// shutting down server with timeout
	log.Println("Application is shutting down")
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Printf("Failed to shutting donw server. Error: %v\n", err)
	}
	log.Printf("Server shutdown gracefully\n")
	return nil
}
