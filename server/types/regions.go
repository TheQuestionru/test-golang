package types

import (
	"database/sql/driver"
	"fmt"
	"golang.org/x/text/language"
	"reflect"
)

type LanguageType string

const (
	LanguageTypeInvalid LanguageType = ""
	LanguageRU          LanguageType = "ru"
	LanguageDE          LanguageType = "de"
	LanguageEN          LanguageType = "en"
	LanguageGB          LanguageType = "gb"
)

var (
	ServerLangs = []language.Tag{
		language.Russian,
		language.English,
		language.German,
	}
)

func (e LanguageType) String() string {
	return string(e)
}

func (e LanguageType) Validate() error {
	if e.Clean() == LanguageTypeInvalid {
		return ErrLangInvalid
	}
	return nil
}

func (e LanguageType) Clean() LanguageType {
	switch e {
	case LanguageRU, LanguageDE, LanguageEN:
		return e
	default:
		return LanguageTypeInvalid
	}

	return e
}

// Scan implements the sql.Scanner interface.
func (e *LanguageType) Scan(value interface{}) error {
	switch v := value.(type) {
	default:
		return fmt.Errorf("Unsupported RegionType type, type=%v", reflect.TypeOf(value))
	case nil:
		*e = LanguageTypeInvalid
	case string:
		*e = LanguageType(v)
	case []byte:
		*e = LanguageType(string(v))
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (e LanguageType) Value() (driver.Value, error) {
	return string(e), nil
}
