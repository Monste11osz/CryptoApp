package config

// DbConf содержит настройки для подключения к PostgreSQL.
type DbConf struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Name string `json:"name"`
}

// ApiExchange содержит конфигурацию для работы с внешними API.
//
// Поля:
//   - Url: адрес API для получения цены валюты.
//   - UrlListCoins: адрес API для получения списка доступных монет
type ApiExchange struct {
	Url          string `json:"URL_PRICE"`
	UrlListCoins string `json:"URL_LIST_COIN"`
}
