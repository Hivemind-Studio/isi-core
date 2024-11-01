package enum

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type Gender int

const (
	Man   Gender = 1
	Woman Gender = 2
)

func (g Gender) IsValid() bool {
	switch g {
	case Man, Woman:
		return true
	}
	return false
}

func (g Gender) Value() (driver.Value, error) {
	if !g.IsValid() {
		return nil, fmt.Errorf("invalid Gender value: %d", g)
	}
	return int(g), nil
}

func (g *Gender) Scan(value interface{}) error {
	val, ok := value.(int64)
	if !ok {
		return errors.New("invalid data type for Gender")
	}
	*g = Gender(val)
	if !g.IsValid() {
		return fmt.Errorf("invalid Gender value: %d", val)
	}
	return nil
}
