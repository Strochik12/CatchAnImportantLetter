package models

import (
	"strings"
	"time"
)

type Email struct {
	ID        ID
	MessageID string            // ID письма из IMAP
	From      string            // Отправитель
	To        []string          // Получатели
	Subject   string            // Тема письма
	Body      string            // Текстовое поле
	HTML      string            // HTML тело
	Date      time.Time         // Дата получения
	Headers   map[string]string // Заголовки
	Links     []string          // Ссылки из письма
	Size      int               // Размер в байтах
	Read      bool              // Прочитано ли
}

// NewEmail создает новый Email с предзаполнеными полями
func NewEmail() *Email {
	return &Email{
		ID:      GenerateID(),
		To:      make([]string, 0),
		Headers: make(map[string]string),
		Links:   make([]string, 0),
		Date:    time.Now(),
	}
}

// ExtractDomain извлекает домен отправителя
func (e *Email) ExtractDomain() string {
	// Разбиваем email на части и берем домен
	// example@hse.ru -> hse.ru
	parts := strings.Split(e.From, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// HasLink проверяет наличие ссылок
func (e *Email) HasLink(pattern string) bool {
	for _, link := range e.Links {
		if strings.Contains(link, pattern) {
			return true
		}
	}
	return false
}
