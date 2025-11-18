package config

import (
	"time"

	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
)

// Config - основная структура конфигурации
type Config struct {
	IMAP       IMAPConfig       `yaml:"imap"`
	Rules      []*models.Rule   `yaml:"rules"`
	Logging    LoggingConfig    `yaml:"logging,omitempty"`
	Monitoring MonitoringConfig `yaml:"monitoring,omitempty"`
}

// IMAPConfig - настройки почтового сервера
type IMAPConfig struct {
	Server         string `yaml:"server"`
	Port           int    `yaml:"port,omitempty"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	Mailbox        string `yaml:"mailbox,omitempty"`
	TLS            bool   `yaml:"tls,omitempty"`
	TimeoutSeconds int    `yaml:"timeout_seconds,omitempty"`
}

// LoggingConfig - настройки логирования
type LoggingConfig struct {
	Level    string `yaml:"level"`  // "debug", "info", "warn", "error"
	Format   string `yaml:"format"` // "json", "text"
	FilePath string `yaml:"file_path,omitempty"`
}

// MonitoringConfig - настройки мониторинга
type MonitoringConfig struct {
	CheckIntervalSeconds int `yaml:"check_interval_seconds,omitempty"`
	MaxEmails            int `yaml:"max_emails,omitempty"`
	RetryAttempts        int `yaml:"retry_attempts,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		IMAP: IMAPConfig{
			Server:         "imap.yandex.ru", // smtp.yandex.ru
			Port:           993,              // 465
			Mailbox:        "INBOX",
			TLS:            true,
			TimeoutSeconds: 30,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
		Monitoring: MonitoringConfig{
			CheckIntervalSeconds: 30,
			MaxEmails:            20,
			RetryAttempts:        3,
		},
	}
}

// Вспомогательный метод для получения duration
func (m *MonitoringConfig) GetCheckInterval() time.Duration {
	return time.Duration(m.CheckIntervalSeconds) * time.Second
}
