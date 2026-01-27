package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// StringSlice stores []string as JSON in MySQL (JSON/TEXT).
// It implements sql.Scanner and driver.Valuer for GORM.
type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	// store nil as empty array (not null)
	if s == nil {
		s = StringSlice{}
	}
	b, err := json.Marshal([]string(s))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (s *StringSlice) Scan(value interface{}) error {
	if s == nil {
		return fmt.Errorf("StringSlice: Scan on nil receiver")
	}
	switch v := value.(type) {
	case nil:
		*s = StringSlice{}
		return nil
	case []byte:
		if len(v) == 0 {
			*s = StringSlice{}
			return nil
		}
		return json.Unmarshal(v, s)
	case string:
		if v == "" {
			*s = StringSlice{}
			return nil
		}
		return json.Unmarshal([]byte(v), s)
	default:
		return fmt.Errorf("StringSlice: unsupported Scan type %T", value)
	}
}

