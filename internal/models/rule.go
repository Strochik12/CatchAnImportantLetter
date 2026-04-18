package models

type Rule struct {
	ID         ID           `yaml:"id" json:"id"`
	Name       string       `yaml:"name" json:"name"`
	Enabled    bool         `yaml:"enabled" json:"enabled"`
	Conditions []Condition  `yaml:"conditions" json:"conditions"`
	Actions    []ActionType `yaml:"actions" json:"actions"`
	Priority   int          `yaml:"priority" json:"priority"`
	MinScore   int          `yaml:"min_score" json:"min_score"`
}

// Condition - условие для правила
type Condition struct {
	Type     ConditionType `yaml:"type" json:"type"`
	Field    string        `yaml:"field,omitempty" json:"field,omitempty"`
	Operator Operator      `yaml:"operator" json:"operator"`
	Value    string        `yaml:"value" json:"value"`
	Weight   int           `yaml:"weight" json:"weight"`
}

// NewRule создает новое правило с предзаполнеными полями
func NewRule(name string) *Rule {
	return &Rule{
		ID:         GenerateID(),
		Name:       name,
		Enabled:    true,
		Conditions: make([]Condition, 0),
		Actions:    make([]ActionType, 0),
		Priority:   50,
		MinScore:   60,
	}
}
