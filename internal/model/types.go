package model

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type StringSlice []string

func (s *StringSlice) Scan(src interface{}) error {
	switch src := src.(type) {
	case string:
		*s = strings.Split(src, ",")
		return nil
	case []byte:
		*s = strings.Split(string(src), ",")
		return nil
	default:
		return errors.New("unsupported type")
	}
}

func (s StringSlice) Value() (driver.Value, error) {
	return strings.Join(s, ","), nil
}
