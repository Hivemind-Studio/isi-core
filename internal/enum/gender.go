package enum

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type Gender int

const (
	Male   Gender = 1
	Female Gender = 2
)

func (g Gender) IsValid() bool {
	switch g {
	case Male, Female:
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
