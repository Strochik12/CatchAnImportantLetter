package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

// Load загружает конфигурацию из файла
func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = findConfigPath()
		if configPath == "" {
			return nil, fmt.Errorf("config file not found")
		}
	}

	// Читаем файл
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Создаем конфиг с значениями по умолчанию
	cfg := DefaultConfig()

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Валидация
	if err := Validate(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// findConfigPath ищет конфиг в стандартных местах
func findConfigPath() string {
	possiblePaths := []string{
		"./configs/config.yaml",
		"./config.yaml",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "" // конфиг не найден
}

// setDefaults устанавливает значения по умолчанию
func setDefaults(v *viper.Viper) {
	// IMAP defaults
	v.SetDefault("imap.server", "outlook.office365.com")
	v.SetDefault("imap.port", 993)
	v.SetDefault("imap.mailbox", "INBOX")
	v.SetDefault("imap.tls", true)
	v.SetDefault("imap.timeout", "30s")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "text")

	// Monitoring defaults
	v.SetDefault("monitoring.check_interval_seconds", 30)
	v.SetDefault("monitoring.max_emails", 100)
	v.SetDefault("monitoring.retry_attempts", 3)
}
