package models

import "time"

type Rule struct {
	ID          ID          `yaml:"id" json:"id"`
	Name        string      `yaml:"name" json:"name"`
	Description string      `yaml:"description" json:"description"`
	Enabled     bool        `yaml:"enabled" json:"enabled"`
	Conditions  []Condition `yaml:"conditions" json:"conditions"`
	Actions     []Action    `yaml:"actions" json:"actions"`
	Priority    int         `yaml:"priority" json:"priority"`
	MinScore    int         `yaml:"min_score" json:"min_score"`
	CreatedAt   time.Time   `yaml:"created_at" json:"created_at"`
	UpdatedAt   time.Time   `yaml:"updated_at" json:"updated_at"`
}

// Condition - условие для правила
type Condition struct {
	Type     ConditionType `yaml:"type" json:"type"`
	Field    string        `yaml:"field,omitempty" json:"field,omitempty"`
	Operator Operator      `yaml:"operator" json:"operator"`
	Value    string        `yaml:"value" json:"value"`
	Weight   int           `yaml:"weight" json:"weight"`
}

// Action - действие при срабатывания правила
type Action struct {
	Type   ActionType    `yaml:"type" json:"type"`
	Config ActionConfig  `yaml:"config" json:"config"`
	Delay  time.Duration `yaml:"delay,omitempty" json:"delay,omitempty"`
}

// ActionConfig - конфигурация действий
type ActionConfig struct {
	Telegram *TelegramActionConfig `yaml:"telegram,omitempty" json:"telegram,omitempty"`
	SMS      *SMSActionConfig      `yaml:"sms,omitempty" json:"sms,omitempty"`
	Webhook  *WebhookActionConfig  `yaml:"webhook,omitempty" json:"webhook,omitempty"`
}

type TelegramActionConfig struct {
	ChatID  int64  `yaml:"chat_id" json:"chat_id"`
	Message string `yaml:"message,omitempty" json:"message,omitempty"`
}

type SMSActionConfig struct {
	Phone   string `yaml:"phone" json:"phone"`
	Message string `yaml:"message,omitempty" json:"message,omitempty"`
}

type WebhookActionConfig struct {
	URL     string            `yaml:"url" json:"url"`
	Method  string            `yaml:"method,omitempty" json:"method,omitempty"` // GET, POST, PUT
	Headers map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Body    string            `yaml:"body,omitempty" json:"body,omitempty"`
	Timeout time.Duration     `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

// NewRule создает новое правило с предзаполнеными полями
func NewRule(name string) *Rule {
	now := time.Now()
	return &Rule{
		ID:         GenerateID(),
		Name:       name,
		Enabled:    true,
		Conditions: make([]Condition, 0),
		Actions:    make([]Action, 0),
		Priority:   50,
		MinScore:   60,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// AddCondition добавляет условие к правилу
func (r *Rule) AddCondition(condition Condition) {
	r.Conditions = append(r.Conditions, condition)
	r.UpdatedAt = time.Now()
}

// AddAction добавляет действие к правилу
func (r *Rule) AddAction(action Action) {
	r.Actions = append(r.Actions, action)
	r.UpdatedAt = time.Now()
}
