package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testYTask/internal/common"
	"testYTask/internal/domain/interfaces"
	"testYTask/internal/domain/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MajorHandler struct {
	majorRepository interfaces.MajorRepositoryI
	mapCoins        map[string]bool
}

func NewMajorHandler(majorRepository interfaces.MajorRepositoryI, mapCoins map[string]bool) *MajorHandler {
	return &MajorHandler{
		majorRepository: majorRepository,
		mapCoins:        mapCoins,
	}
}

// AddingCoin обрабатывает запрос на добавление новой монеты.
//
// Маршрут: POST /api/v1/currency/add
//
// Параметры запроса (JSON):
//   - name_coin: название монеты (string)
//
// Возможные ответы:
//   - 200 OK: монета успешно добавлена.
//   - 400 Bad Request: некорректные входные данные.
//   - 404 Not Found: монета не найдена.
//   - 500 Internal Server Error: ошибка сервиса при добавлении.
func (h *MajorHandler) AddingCoin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), common.ReqTimeout)
	defer cancel()

	zap.L().Info("Start adding coins...")

	coin := new(models.Coin)

	if err := c.ShouldBind(&coin); err != nil || coin.NameCoin == common.Empty {
		zap.L().Error("ShouldBind error", zap.Error(err))
		common.ResponseBadRequest(c, "Incorrect input data")
		return
	}

	if !h.mapCoins[strings.ToLower(coin.NameCoin)] {
		zap.L().Error("Map coin not found", zap.String("name:", coin.NameCoin))
		common.ResponseBadRequest(c, "Coin not found")
		return
	}

	if err := h.majorRepository.AddingNewCoin(ctx, coin); err != nil {
		zap.L().Error("AddingNewCoin error", zap.Error(err), zap.String("name:", coin.NameCoin))
		common.ResponseServerError(c, "Service error while adding")
		return
	}
	common.ResponseSuccess(c, fmt.Sprintf("Coin '%s' added", coin.NameCoin), struct{}{})

	zap.L().Info("Successful coin addition")
}

// DeleteCoin обрабатывает запрос на удаление монеты.
//
// Маршрут: DELETE /api/v1/currency/remove
//
// Параметры запроса (JSON):
//   - name_coin: название монеты (string)
//
// Возможные ответы:
//   - 200 OK: монета успешно удалена.
//   - 400 Bad Request: некорректные входные данные.
//   - 404 Not Found: монета отсутствует в списке.
//   - 500 Internal Server Error: ошибка сервиса при удалении.
func (h *MajorHandler) DeleteCoin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), common.ReqTimeout)
	defer cancel()

	zap.L().Info("Start deletion coins...")

	coin := new(models.Coin)

	if err := c.ShouldBind(&coin); err != nil {
		zap.L().Error("ShouldBind error", zap.Error(err))
		common.ResponseBadRequest(c, "Incorrect input data")
		return
	}

	if err := h.majorRepository.DeleteCoin(ctx, coin); err != nil {
		switch {
		case errors.Is(err, common.ErrCoinNotFound):
			common.ResponseBadRequest(c, "This coin is not on the list")
			return
		default:
			zap.L().Error("DeleteCoin error", zap.Error(err), zap.String("name", coin.NameCoin))
			common.ResponseServerError(c, "Service error while deleting")
			return
		}
	}
	common.ResponseSuccess(c, "Coin deleted", struct{}{})

	zap.L().Info("Successful coin deletion")
}

// GetPriceForCoin обрабатывает запрос на получение цены указанной монеты.
//
// Маршрут: POST /api/v1/currency/price
//
// Параметры запроса (JSON):
//   - coin: название монеты (string)
//   - timestamp: метка времени (int)
//
// Возможные ответы:
//   - 200 OK: цена успешно получена.
//   - 400 Bad Request: некорректные входные данные.
//   - 404 Not Found: цена не найдена.
//   - 500 Internal Server Error: ошибка сервиса при получении данных.
func (h *MajorHandler) GetPriceForCoin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), common.ReqTimeout)
	defer cancel()

	zap.L().Info("Start getting price for coin...")

	req := new(models.PriceRequest)

	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("ShouldBind error", zap.Error(err), zap.String("coin:", req.Coin), zap.Int64("Time:", req.Timestamp))
		common.ResponseBadRequest(c, "Invalid request")
		return
	}

	if req.Coin == common.Empty || req.Timestamp <= common.Zero {
		zap.L().Error("GetPriceForCoin invalid coin or timestamp", zap.String("coin:", req.Coin), zap.Int64("Time:", req.Timestamp))
		common.ResponseBadRequest(c, "Required fields: coin and timestamp")
		return
	}

	data, err := h.majorRepository.GetPrice(ctx, req)
	if err != nil {
		if errors.Is(err, common.ErrPriceNotFound) {
			zap.L().Info("GetPriceForCoin not found", zap.String("coin:", req.Coin))
			common.ResponseNotFound(c, "Price not found")
			return
		} else {
			zap.L().Error("DB error", zap.Error(err))
			common.ResponseServerError(c, "Error while receiving data")
			return
		}
	}
	common.ResponseSuccess(c, common.Empty, data)

	zap.L().Info("Successful coin getPrice")
}
