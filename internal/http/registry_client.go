package http

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testYTask/internal/common"
	"testYTask/internal/config"
	"testYTask/internal/domain/models"

	"go.uber.org/zap"
)

// RegistryClient отвечает за отправку HTTP-запросов к внешним сервисам
type RegistryClient struct {
	httpClient *http.Client
	cfg        *config.ApiExchange
}

// NewRegistryClient создаёт и инициализирует новый экземпляр RegistryClient с заданной конфигурацией.
func NewRegistryClient(cfg *config.ApiExchange) *RegistryClient {
	return &RegistryClient{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: common.RequestTimeout,
		},
		cfg: cfg,
	}
}

// LoadValidCoins загружает список всех валидных монет с внешнего API.
//
// Параметры:
//   - ctx: контекст запроса для управления временем выполнения и отменой.
//
// Возвращает:
//   - map[string]bool: словарь. Ключ - название монеты, значение - true.
//   - error: ошибка при получении или обработке данных.
func (r *RegistryClient) LoadValidCoins(ctx context.Context) (map[string]bool, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, common.ReqTimeValidCoins)
	defer cancel()

	var coins []struct {
		ID     string `json:"id"`
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
	}

	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodGet, r.cfg.UrlListCoins, nil)
	if err != nil {
		zap.L().Error("Error creating HTTP request", zap.Error(err))
		return nil, err
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		zap.L().Error("Error executing HTTP request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("Error reading HTTP response body", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("Unexpected status code: %d", resp.StatusCode)

		zap.L().Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err = json.Unmarshal(body, &coins); err != nil {
		zap.L().Error("Error unmarshalling HTTP response body", zap.Error(err))
		return nil, err
	}

	validCoins := make(map[string]bool)

	for _, c := range coins {
		validCoins[c.ID] = true
	}
	return validCoins, nil
}

// CurrentData получает текущие данные о цене указанной монеты с внешнего API.
//
// Параметры:
//   - ctx: контекст запроса для управления временем выполнения и отменой.
//   - coin: название монеты (например, "bitcoin").
//
// Возвращает:
//   - *models.CoinPrice: структура с ценой монеты, валютой и временем получения.
//   - error: ошибка при получении или обработке данных.
func (r *RegistryClient) CurrentData(ctx context.Context, coin string) (*models.CoinPrice, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, common.ReqTimePrice)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodGet, fmt.Sprintf(r.cfg.Url, coin), nil)
	if err != nil {
		zap.L().Error("Error creating HTTP request", zap.Error(err))
		return nil, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		zap.L().Error("Error executing HTTP request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("Error reading HTTP response body", zap.Error(err))
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("Unexpected status code: %d", resp.StatusCode)

		zap.L().Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	priceResponse := make(map[string]*models.CoinPrice)

	if err = json.Unmarshal(body, &priceResponse); err != nil {
		zap.L().Error("Error unmarshalling HTTP response body", zap.Error(err))
		return nil, err
	}
	return priceResponse[coin], nil
}
