package config

import "time"

// AppConf содержит основные настройки приложения.
//
// Поля:
//   - Mode: режим работы приложения (например, "debug" или "release").
//   - Code: уникальный код приложения.
//   - Port: порт, на котором запускается сервер.
//   - Stage: стадия запуска приложения (например, "production", "staging").
//   - RTO: время ожидания чтения запроса (Read Timeout).
//   - WTO: время ожидания записи ответа (Write Timeout).
type AppConf struct {
	Mode  string        `json:"APP_MODE"`
	Code  string        `json:"APP_CODE"`
	Port  int           `json:"APP_PORT"`
	Stage string        `json:"APP_STAGE"`
	RTO   time.Duration `json:"APP_RTO"`
	WTO   time.Duration `json:"APP_WTO"`
}
