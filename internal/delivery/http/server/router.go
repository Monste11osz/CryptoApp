package server

import (
	"testYTask/docs"
	"testYTask/internal/config"
	"testYTask/internal/delivery/http"
	"testYTask/internal/delivery/http/v1/handlers"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// Navigator управляет маршрутизацией HTTP запросов.
type Navigator struct {
	cfg    *config.Config
	Engine *gin.Engine
}

// NewNavigator создаёт новый Navigator и настраивает движок Gin.
//
// Параметры:
//   - cfg: конфигурация HTTP сервера
//
// Возвращает указатель на Navigator.
func NewNavigator(cfg *config.Config) *Navigator {
	gin.SetMode(cfg.App.Mode)

	docs.SwaggerInfo.Title = cfg.Swagger.Title
	docs.SwaggerInfo.Description = cfg.Swagger.Description
	docs.SwaggerInfo.Version = cfg.Swagger.Version
	docs.SwaggerInfo.BasePath = cfg.Swagger.BaseURL

	Engine := gin.New()
	Engine.Use(
		ginzap.GinzapWithConfig(
			zap.L(),
			&ginzap.Config{
				TimeFormat: time.RFC1123Z,
				UTC:        true,
				SkipPaths:  []string{"/healthz"},
			},
		),
	)
	Engine.Use(ginzap.RecoveryWithZap(zap.L(), true))
	Engine.Use(
		cors.New(cors.Config{
			AllowOrigins: cfg.Cors.AllowOrigins,
			AllowMethods: cfg.Cors.AllowMethods,
			AllowHeaders: cfg.Cors.AllowHeaders,
		}),
	)
	Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Navigator{
		cfg:    cfg,
		Engine: Engine,
	}
}

// RegisterRoutes регистрирует маршруты для HTTP-сервера.
//
// Параметры:
//   - commonHandler: обработчик для общих маршрутов (например, проверки состояния)
//   - majorHandler: обработчик для основных бизнес-операций
func (n *Navigator) RegisterRoutes(commonHandler *http.CommonHandler, majorHandler *handlers.MajorHandler) {
	n.Engine.GET("healthz", commonHandler.HealthCheck)
	api := n.Engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			currency := v1.Group("/currency")
			{
				currency.POST("/price", majorHandler.GetPriceForCoin)
				currency.POST("/add", majorHandler.AddingCoin)
				currency.DELETE("/remove", majorHandler.DeleteCoin)
			}
		}
	}
}
