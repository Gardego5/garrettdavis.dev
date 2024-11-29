//go:generate msgp
package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time struct{ time.Time }

func (t *Time) Scan(value any) error {
	if str, ok := value.(string); !ok {
		return fmt.Errorf("expected string as internal db value, got %T", value)
	} else {
		return t.Time.UnmarshalText([]byte(str))
	}
}

func (t Time) Value() (driver.Value, error) {
	if time, err := t.Time.MarshalText(); err != nil {
		return nil, err
	} else {
		return string(time), nil
	}
}
