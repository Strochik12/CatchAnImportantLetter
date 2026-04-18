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
		if rule == nil {
			return fmt.Errorf("rule %d is nil", i)
		}

		// Валидируем само правило
		if err := validateRule(rule); err != nil {
			return fmt.Errorf("rule '%s' validation failed: %w", rule.Name, err)
		}

		// Проверяем уникальность имен правил
		if ruleNames[rule.Name] {
			return fmt.Errorf("duplicate rule name: %s", rule.Name)
		}
		ruleNames[rule.Name] = true
	}

	return nil
}

// Validate проверяет правило на валидность
func validateRule(r *models.Rule) error {
	if r.Name == "" {
		return fmt.Errorf("rule name cannot be empty")
	}
	if len(r.Conditions) == 0 {
		return fmt.Errorf("rule must have at least one condition")
	}
	if len(r.Actions) == 0 {
		return fmt.Errorf("rule must have at least one action")
	}
	if r.MinScore < 0 || r.MinScore > 100 {
		return fmt.Errorf("min_score must be between 0 and 100")
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
