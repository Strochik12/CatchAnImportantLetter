package models

import "testing"

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		name           string
		email          *Email
		expectedDomain string
	}{
		{
			name:           "Домен есть",
			email:          &Email{From: "example@hse.ru"},
			expectedDomain: "hse.ru",
		},
		{
			name:           "Домена нет",
			email:          &Email{From: "example"},
			expectedDomain: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultDomain := tt.email.ExtractDomain()
			if resultDomain != tt.expectedDomain {
				t.Errorf("incorrect result, expected: '%v', got: '%v'", tt.expectedDomain, resultDomain)
			}
		})
	}
}

func TestHasLink(t *testing.T) {
	tests := []struct {
		name           string
		email          *Email
		link           string
		expectedResult bool
	}{
		{
			name:           "Есть нужная ссылка",
			email:          &Email{Links: []string{"github.com/Strochik12", "https://open.spotify.com", "youtube.com"}},
			link:           "spotify.com",
			expectedResult: true,
		},
		{
			name:           "Нет нужной ссылки",
			email:          &Email{Links: []string{"github.com/Strochik12", "https://open.spotify.com", "youtube.com"}},
			link:           "gmail.com",
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.email.HasLink(tt.link)
			if result != tt.expectedResult {
				t.Errorf("incorrect result")
			}
		})
	}
}
