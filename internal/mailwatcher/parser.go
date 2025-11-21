package mailwatcher

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/quotedprintable"
	"strings"

	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"golang.org/x/text/encoding/charmap"
)

func init() {
	// Регистрируем поддержку разных кодировок
	charset.RegisterEncoding("windows-1251", charmap.Windows1251)
	charset.RegisterEncoding("koi8-r", charmap.KOI8R)
}

// parseMessage преобразует imap.Message в наш models.Email
func parseMessage(msg *imap.Message) (*models.Email, error) {
	if msg.Envelope == nil {
		return nil, fmt.Errorf("письмо не содержит envelope")
	}

	email := models.NewEmail()
	email.MessageID = msg.Envelope.MessageId
	email.Subject = msg.Envelope.Subject
	email.Date = msg.Envelope.Date

	// Обрабатываем отправителя
	if len(msg.Envelope.From) > 0 {
		from := msg.Envelope.From[0]
		email.From = formatAddress(from)
	}

	// Обрабатываем отправителей
	for _, to := range msg.Envelope.To {
		email.To = append(email.To, formatAddress(to))
	}

	// Парсим тело письма
	if err := parseBody(msg, email); err != nil {
		return email, fmt.Errorf("ошибка парсинга тела: %w", err)
	}

	return email, nil
}

// parseBody парсит тело письма
func parseBody(msg *imap.Message, email *models.Email) error {
	// Пробуем разные подходы к получению тела письма

	// 1. Сначала пробуем получить raw body и распарсить как MIMEBody
	if section := msg.GetBody(&imap.BodySectionName{}); section != nil {
		if err := parseMIMEBody(section, email); err == nil {
			return nil // Успешно распарсили через MIME
		}
	}

	// 2. Если MIME не сработал, пробуем простой TEXT
	if section := msg.GetBody(&imap.BodySectionName{
		BodyPartName: imap.BodyPartName{Specifier: imap.TextSpecifier},
	}); section != nil {
		if content, err := io.ReadAll(section); err == nil {
			email.Body = string(content)
			return nil
		}
	}

	// 3. Последняя попытка - пробуем явно получить multipart MIME структуру
	if section := msg.GetBody(&imap.BodySectionName{
		BodyPartName: imap.BodyPartName{Specifier: imap.MIMESpecifier},
	}); section != nil {
		if err := parseMIMEBody(section, email); err == nil {
			return nil // Успешно распарсили через MIME
		}
	}

	return fmt.Errorf("не удалось извлечь тело письма")
}

// parseMIMEBody парсит MIME структуру письма
func parseMIMEBody(section io.Reader, email *models.Email) error {
	mr, err := mail.CreateReader(section)
	if err != nil {
		return err
	}
	defer mr.Close()

	// Читаем заголовки (From, To, Subject, Date и тд)
	for key, values := range mr.Header.Map() {
		email.Headers[key] = strings.Join(values, ", ")
	}

	// Парсим части письма (текст, HTML, вложения)
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Пропускаем битые части
		}

		body, err := io.ReadAll(part.Body)
		if err != nil {
			continue
		}

		contentType := part.Header.Get("Content-Type")
		// Повыводив письма в сыром виде я заметил, что в письмах яндекса используется base64 кодирование
		contentEncoding := part.Header.Get("Content-Transfer-Encoding")

		// Декодируем содержимое
		content, err := decodeContent(body, contentEncoding)
		if err != nil {
			log.Printf("Ошибка декодировки письма: %v", err)
			continue
		}

		// Сохраняем в зависимости от типа
		switch {
		case strings.Contains(contentType, "text/plain"):
			email.Body = content
		case strings.Contains(contentType, "text/html"):
			email.HTML = content
			email.Links = extractLinks(content)
		}
	}

	return nil
}

// extractLinks - извлекает ссылки из HTML (Пока что упрощенная)
func extractLinks(html string) []string {
	var links []string

	// Простой поиск href="..."
	start := 0
	for {
		idx := strings.Index(html[start:], "href=\"")
		if idx == -1 {
			break
		}

		start += idx + 6
		end := strings.Index(html[start:], "\"")
		if end == -1 {
			break
		}

		link := html[start : start+end]
		if strings.HasPrefix(link, "http") {
			links = append(links, link)
		}

		start += end + 1
	}

	return links
}

// formatAddress форматирует адрес email
func formatAddress(addr *imap.Address) string {
	if addr.PersonalName != "" {
		return fmt.Sprintf("%s <%s@%s>", addr.PersonalName, addr.MailboxName, addr.HostName)
	}
	return fmt.Sprintf("%s@%s", addr.MailboxName, addr.HostName)
}

// decodeContent декодирует используемые в Яндекс почте кодировки
func decodeContent(body []byte, encoding string) (string, error) {
	switch strings.ToLower(encoding) {
	case "base64":
		decoded, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			// Если указано base64, но содержимое не base64 - возвращаем как есть
			return string(body), nil
		}
		return string(decoded), nil

	case "quoted-printable":
		decoded, err := io.ReadAll(quotedprintable.NewReader(bytes.NewReader(body)))
		if err != nil {
			return "", fmt.Errorf("ошибка декодирования quoted-printable: %w", err)
		}
		return string(decoded), nil

	case "7bit", "8bit", "binary", "":
		// Без кодировки или простые текстовые кодировки
		return string(body), nil

	default:
		return "", fmt.Errorf("неподдерживаемая кодировка: %s", encoding)
	}
}
