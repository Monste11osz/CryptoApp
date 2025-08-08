package http

import (
	"context"
	"net/http"
	"testYTask/internal/domain/interfaces"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommonHandler обрабатывает запросы общего характера, например, проверку работоспособности сервиса.
type CommonHandler struct {
	sentinels []interface{}
}

// NewCommonHandler создаёт новый экземпляр CommonHandler.
//
// Параметры:
//   - sentinels: список компонентов, поддерживающих проверку работоспособности
//
// Возвращает указатель на CommonHandler.
func NewCommonHandler(sentinels ...interface{}) *CommonHandler {
	return &CommonHandler{
		sentinels: sentinels,
	}
}

// HealthCheck обрабатывает запрос /healthz и проверяет работоспособность компонентов.
//
// Параметры:
//   - c: контекст запроса Gin
func (h *CommonHandler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, sentinel := range h.sentinels {
		if sentinel == nil {
			continue
		}
		switch s := sentinel.(type) {
		case interfaces.PingerI:
			if err := s.Ping(ctx); err != nil {
				zap.L().Error("Database or kafka connection error", zap.Error(err))
				c.AbortWithStatus(http.StatusServiceUnavailable)
				return
			}
		default:
			zap.L().Error("Unknown sentinel type", zap.Any("type", s))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	c.Status(http.StatusOK)
}
