package mailwatcher

import (
	"fmt"
	"log"

	"github.com/Strochik12/CatchAnImportantLetter/internal/config"
	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// Client обертка вокруг IMAP соединения
type Client struct {
	config    *config.Config
	client    *client.Client
	connected bool
	lastUid   uint32
}

// NewIMAP создает новый IMAP клиент
func NewIMAPClient(cfg *config.Config) *Client {
	return &Client{
		config:    cfg,
		connected: false,
		lastUid:   0,
	}
}

// Connect устанавливает соединение с IMAP сервером
func (c *Client) Connect() error {
	var err error
	addr := fmt.Sprintf("%s:%d", c.config.IMAP.Server, c.config.IMAP.Port)

	log.Printf("Подключение к IMAP серверу: %s ...", addr)

	if c.config.IMAP.TLS {
		c.client, err = client.DialTLS(addr, nil)
	} else {
		c.client, err = client.Dial(addr)
	}

	if err != nil {
		return fmt.Errorf("не удалось подключиться к серверу: %w", err)
	}

	if err := c.client.Login(c.config.IMAP.Username, c.config.IMAP.Password); err != nil {
		c.client.Logout()
		return fmt.Errorf("ошибка авторизации: %w", err)
	}

	c.connected = true
	log.Printf("Успешное подключение к почтовому ящику")

	return nil
}

// GetNewEmails возвращает новые письма
func (c *Client) GetNewEmails() ([]*models.Email, error) {
	if !c.connected {
		return nil, fmt.Errorf("клиент не подключен")
	}

	// Выбираем почтовый ящик
	mailbox, err := c.client.Select(c.config.IMAP.Mailbox, false)
	if err != nil {
		return nil, fmt.Errorf("ошибка выбора ящика: %w", err)
	}

	// Если нет новых писем
	if mailbox.Messages == 0 {
		return []*models.Email{}, nil
	}

	// Получаем письма только с UID больше последнего обработанного
	// Создаем команду UidSearch для поиска UID > lastUid
	criteria := &imap.SearchCriteria{
		Uid: new(imap.SeqSet),
	}
	criteria.Uid.AddRange(c.lastUid+1, 0) // От lastUid + 1 до конца

	uids, err := c.client.UidSearch(criteria)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска писем: %w", err)
	}

	// Если нет новых писем
	if len(uids) == 0 {
		return []*models.Email{}, nil
	}

	// Либо берём последние MaxEmails писем
	from := c.lastUid + 1
	if len(uids) > c.config.Monitoring.MaxEmails {
		from = uids[len(uids)-c.config.Monitoring.MaxEmails]
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, 0)

	// Запрашиваем заголовки и тела писем
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchInternalDate, section.FetchItem()}

	go func() {
		done <- c.client.UidFetch(seqset, items, messages)
	}()

	var emails []*models.Email
	for msg := range messages {
		email, err := parseMessage(msg)

		if err != nil {
			log.Printf("Ошибка парсинга письма: %v", err)
			continue
		}

		emails = append(emails, email)
		c.lastUid = msg.Uid
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("ошибка получения писем: %w", err)
	}

	log.Printf("Найдено писем: %d", len(emails))
	return emails, nil
}

// Close закрывает соединение
func (c *Client) Close() error {
	if c.connected {
		c.connected = false
		return c.client.Logout()
	}
	return nil
}

// IsConnected возвращает статус подключения
func (c *Client) IsConnected() bool {
	return c.connected
}
