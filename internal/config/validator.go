package config

import (
	"fmt"

	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
)

// Validate проверяет корректность конфигурации
func Validate(cfg *Config) error {
	if err := validateIMAP(&cfg.IMAP); err != nil {
		return fmt.Errorf("IMAP config error: %w", err)
	}

	if err := validateRules(cfg.Rules); err != nil {
		return fmt.Errorf("rules config error: %w", err)
	}

	if err := validateMonitoring(&cfg.Monitoring); err != nil {
		return fmt.Errorf("monitoring config error: %w", err)
	}

	return nil
}

func validateIMAP(imap *IMAPConfig) error {
	if imap.Server == "" {
		return fmt.Errorf("server is required")
	}
	if imap.Username == "" {
		return fmt.Errorf("username is required")
	}
	if imap.Password == "" {
		return fmt.Errorf("password is required")
	}
	if imap.Port <= 0 || imap.Port > 65535 {
		return fmt.Errorf("invalid port: %d", imap.Port)
	}
	return nil
}

func validateRules(rules []*models.Rule) error {
	if len(rules) == 0 {
		return fmt.Errorf("at least one rule is required")
	}

	ruleNames := make(map[string]bool)
	for i, rule := range rules {
		if rule.Name == "" {
			return fmt.Errorf("rule %d: name is required", i)
		}

		// Проверяем уникальность имен правил
		if ruleNames[rule.Name] {
			return fmt.Errorf("duplicate rule name: %s", rule.Name)
		}
		ruleNames[rule.Name] = true

		// Валидируем само правило
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("rule '%s' validation failed: %w", rule.Name, err)
		}
	}

	return nil
}

func validateMonitoring(monitoring *MonitoringConfig) error {
	if monitoring.CheckIntervalSeconds < 5 {
		return fmt.Errorf("check_interval_seconds too small: %v", monitoring.CheckIntervalSeconds)
	}
	if monitoring.MaxEmails <= 0 {
		return fmt.Errorf("max_emails must be positive")
	}
	return nil
}
