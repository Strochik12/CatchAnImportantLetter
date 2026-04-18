package config

import (
	"testing"

	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
)

func TestValidate(t *testing.T) {
	goodCfg := DefaultConfig()
	goodCfg.IMAP.Username = "qwerty@123.abc"
	goodCfg.IMAP.Password = "qwerty123"
	goodCfg.Rules = []*models.Rule{
		{
			ID:       "urgent",
			Name:     "Срочное",
			Enabled:  true,
			Priority: 100,
			MinScore: 10,
			Conditions: []models.Condition{
				{
					Type:     "subject",
					Operator: "contains",
					Value:    "срочно",
					Weight:   10,
				},
			},
			Actions: []models.ActionType{"telegram"},
		},
	}
	tests := []struct {
		name    string
		wantErr bool
		cfg     Config
	}{
		{
			name:    "Правильный конфиг",
			wantErr: false,
			cfg:     *goodCfg,
		},
		{
			name:    "Нет конфига",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(&tt.cfg)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
