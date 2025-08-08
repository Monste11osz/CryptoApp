package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testYTask/internal/config"
	"testYTask/internal/domain/interfaces"
	"time"

	"github.com/gin-gonic/gin"
)

// nexus реализует интерфейс NexusI и управляет HTTP сервером.
type nexus struct {
	cfg    *config.Config
	server *http.Server
}

// NewNexus создаёт новый экземпляр HTTP сервера.
//
// Параметры:
//   - cfg: конфигурация HTTP сервера
//   - engine: Gin движок с настроенными маршрутами
//
// Возвращает объект, реализующий интерфейс NexusI.
func NewNexus(cfg *config.Config, engine *gin.Engine) interfaces.NexusI {
	return &nexus{
		cfg: cfg,
		server: &http.Server{
			Addr:           fmt.Sprintf(":%d", cfg.App.Port),
			Handler:        engine,
			MaxHeaderBytes: 1 << 20,
			ReadTimeout:    cfg.App.RTO * time.Second,
			WriteTimeout:   cfg.App.WTO * time.Second,
		},
	}
}

// Start запускает HTTP сервер и слушает входящие запросы.
//
// Возвращает ошибку, если сервер завершился с ошибкой.
func (n *nexus) Start() error {
	err := n.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

// Stop корректно останавливает HTTP сервер.
//
// Возвращает ошибку, если остановка сервера завершилась неудачно.
func (n *nexus) Stop() error {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return n.server.Shutdown(ctxWithTimeout)
}
