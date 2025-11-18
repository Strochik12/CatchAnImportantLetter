package config

import (
	"fmt"
	"os"

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
