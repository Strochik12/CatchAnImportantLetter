package mailwatcher

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	stateFile string
}

// NewIMAP создает новый IMAP клиент
func NewIMAPClient(cfg *config.Config) *Client {
	client := &Client{
		config:    cfg,
		connected: false,
		stateFile: "data/mail_state.json",
	}
	client.loadState()
	return client
}

// loadState загружает lastUid из файла сохранения
func (c *Client) loadState() error {
	data, err := os.ReadFile(c.stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			c.lastUid = 0 // Первый запуск
			return nil
		}
		return fmt.Errorf("ошибка чтения файла состояния: %w", err)
	}

	var state struct {
		LastUid uint32 `json:"last_uid"`
	}

	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("ошибка демаршалинга состояния: %w", err)
	}

	c.lastUid = state.LastUid
	return nil
}

// saveState атомарно записывает lastUid в файл сохранения
func (c *Client) saveState() error {
	state := struct {
		LastUid uint32 `json:"last_uid"`
	}{
		LastUid: c.lastUid,
	}

	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("ошибка в маршалинга: %w", err)
	}

	// Создаем директорию, если она не существует
	dir := filepath.Dir(c.stateFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("ошибка создание директории: %w", err)
	}

	// Создаем временный файл для атомарной записи
	tmpFile := c.stateFile + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("ошибка записи временного файла: %w", err)
	}

	// Атомарно заменяем старый файл новым
	if err := os.Rename(tmpFile, c.stateFile); err != nil {
		return fmt.Errorf("ошибка переименовывания файла: %w", err)
	}

	return nil
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

	// Если нет писем вообще
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
	if len(uids) == 0 || (len(uids) == 1 && uids[0] == c.lastUid) {
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
		if msg.Uid > c.lastUid {
			c.lastUid = msg.Uid
		}
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("ошибка получения писем: %w", err)
	}

	if err := c.saveState(); err != nil {
		return nil, fmt.Errorf("ошибка сохранения состояния: %w", err)
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
