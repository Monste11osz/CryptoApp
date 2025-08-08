package config

// CorsConf описывает настройки CORS для приложения.
type CorsConf struct {
	AllowOrigins []string `json:"ALLOW_ORIGINS"`
	AllowMethods []string `json:"ALLOW_METHODS"`
	AllowHeaders []string `json:"ALLOW_HEADERS"`
}

// SwagConf описывает настройки SWAGGER для приложения.
type SwagConf struct {
	Title       string `json:"SWAG_TITLE"`
	Description string `json:"SWAG_DESCRIPTION"`
	Version     string `json:"VERSION"`
	BaseURL     string `json:"BASE_URL"`
}
