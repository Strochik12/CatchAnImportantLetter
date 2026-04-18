package config

import (
	"fmt"
	"testing"
)

type MockFileManager struct{}

func (m MockFileManager) ReadFile(path string) ([]byte, error) {
	switch path {
	case "non-valid.yaml":
		return []byte("qwerty"), nil
	default:
		return nil, fmt.Errorf("file %v not found", path)
	}
}

func (m MockFileManager) CheckFile(path string) bool {
	return path == "./config.yaml"
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		path    string
	}{
		{
			name:    "Путь к несуществующему файлу",
			wantErr: true,
			path:    "non-existant",
		},
		{
			name:    "Невалидный .yaml файл",
			wantErr: true,
			path:    "non-valid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadWithFileManager(tt.path, MockFileManager{})
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}
