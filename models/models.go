package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Date struct {
	time.Time
}

type UpdateMovie struct {
	Title       *string `json:"title"`
	ReleaseYear *int    `json:"releaseYear"`
	Duration    *int    `json:"duration"`
}

const dateLayout = "2006-01-02"

func (d *Date) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	if str == "" {
		return nil
	}
	t, err := time.Parse(dateLayout, str)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format(dateLayout) + `"`), nil
}

func (d Date) Value() (driver.Value, error) {
	return d.Time.Format("2006-01-02"), nil
}

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	default:
		return fmt.Errorf("Cannot scan type %T into date", value)
	}
}
