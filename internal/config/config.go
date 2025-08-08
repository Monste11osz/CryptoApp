package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var (
	once sync.Once
	Conf = &Config{}
)

type Config struct {
	App      *AppConf     `json:"APP"`
	Cors     *CorsConf    `json:"CORS"`
	Db       *DbConf      `json:"DB"`
	Exchange *ApiExchange `json:"EXCHANGE"`
	Swagger  *SwagConf    `json:"SWAGGER"`
}

// GetConfig загружает и возвращает конфигурацию из указанного файла.
//
// Параметры:
//   - path: путь к файлу конфигурации
//
// Возвращает указатель на объект Config.
func GetConfig() *Config {
	once.Do(
		func() {
			Conf = new(Config)
			LoadLocalConf(Conf)
		},
	)
	return Conf
}

// LoadLocalConf загружает конфигурацию из файла по указанному пути.
func LoadLocalConf(Conf *Config) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "./config/conf.json"
	}
	LoadConfig(Conf, path)
}

// LoadConfig загружает конфигурацию приложения из JSON-файла.
//
// Паника:
//   - Если файл не найден или недоступен для чтения.
//   - Если содержимое файла не соответствует ожидаемой структуре Config.
func LoadConfig(Conf *Config, path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("[LoadConfig]: Error| %s", err.Error()))
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(Conf); err != nil {
		panic(fmt.Sprintf("[LoadConfig]: Error| %s", err.Error()))
	}
}
