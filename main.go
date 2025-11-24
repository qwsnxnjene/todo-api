package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func initConfig() {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")

	viper.AutomaticEnv()
}

func main() {
	initConfig()
	initLogger()
	defer Logger.Sync()

	storage := NewStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", storage.ListHandler)
	mux.HandleFunc("POST /tasks", storage.CreateHandler)

	server := &http.Server{
		Addr:    ":" + viper.GetString("PORT"),
		Handler: mux,
	}

	// горутина с запуском сервера
	go func() {
		Logger.Info("Сервер запущен", zap.String("port", viper.GetString("PORT")))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger.Fatal("ошибка сервера", zap.Error(err))
		}
	}()

	// ждём сигнал завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	Logger.Info("Выключаемся...")

	// даём 10 секунд на завершение запросов
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		Logger.Error("сервер принудительно выключен", zap.Error(err))
	} else {
		Logger.Info("сервер красиво завершился")
	}
}
