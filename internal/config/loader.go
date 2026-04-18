package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

type FileManager interface {
	ReadFile(path string) ([]byte, error)
	CheckFile(path string) bool
}

type OSFileManager struct{}

func (m OSFileManager) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (m OSFileManager) CheckFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Load загружает конфигурацию из файла
func Load(configPath string) (*Config, error) {
	return loadWithFileManager(configPath, OSFileManager{})
}

// loadWithFileManager загружает конфигурацию из файла с путём configPath с помощью m FileManager
func loadWithFileManager(configPath string, m FileManager) (*Config, error) {
	if configPath == "" {
		configPath = findConfigPath(m)
		if configPath == "" {
			return nil, fmt.Errorf("config file not found")
		}
	}

	// Читаем файл
	data, err := m.ReadFile(configPath)
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
func findConfigPath(m FileManager) string {
	possiblePaths := []string{
		"./configs/config.yaml",
		"./config.yaml",
	}

	for _, path := range possiblePaths {
		if m.CheckFile(path) {
			return path
		}
	}

	return "" // конфиг не найден
}
