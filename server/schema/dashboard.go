package schema

import (
	"database/sql/driver"
	"fmt"
	"github.com/TheQuestionru/thequestion/server/types"
	"github.com/kapitanov/go-teamcity"
	"github.com/yfronto/newrelic"
	"reflect"
)

// errors

var (
	ErrDashboardElementTypeInvalid      = types.NewError("errors.DashboardElementTypeInvalid")
	ErrDashboardElementMaxCol           = types.NewError("errors.DashboardElementMaxCol")           // "В строке может быть от 0 до 5 элементов"
	ErrDashboardTemplateNoParamsAllowed = types.NewError("erorrs.DashboardTemplateNoParamsAllowed") // "Используйте шаблон отчета без параметров"
)

// structs
type DashboardElementKey struct {
	Type  DashboardElementType `json:"type" db:"type"`
	Title types.NullString     `json:"title" db:"title"`

	ReportTemplateId types.NullInt64 `json:"reportTemplateId" db:"report_template_id"`
}

type DashboardElement struct {
	DashboardElementKey

	Id      int64 `json:"-" id:"true" generated:"true"`
	Deleted bool  `json:"deleted" db:"deleted"`

	CsvLink types.NullString `json:"csvLink" db:"csv_link"`

	RowNumber int32 `json:"rowNumber" db:"row_number"`
	ColNumber int32 `json:"colNumberr" db:"col_number"`

	CreatedAt types.Time     `json:"createdAt" db:"created_at"`
	DeletedAt types.NullTime `json:"deletedAt" db:"deleted_at"`
}

// view

type DashboardView struct {
	Rows []*DashboardRowView `json:"rows"`
}

type DashboardRowView struct {
	Elements []*DashboardElementView `json:"elements"`
}

type DashboardElementView struct {
	*DashboardElement

	Realtime types.NullInt64   `json:"realtime"`
	Servers  []newrelic.Server `json:"servers,omitempty"`
	Builds   []teamcity.Build  `json:"builds,omitempty"`
}

// forms

type DashboardElementsForm struct {
	Rows []DashboardRow `json:"rows"`
}

type DashboardRow struct {
	Elements []DashboardElementKey `json:"elements"`
}

// types

type DashboardElementType string

const (
	DashboardElementTypeInvalid        DashboardElementType = ""
	DashboardElementTypeReportTemplate DashboardElementType = "report"
	DashboardElementTypeGARealtime     DashboardElementType = "ga-realtime"
	DashboardElementTypeNRServers      DashboardElementType = "nr-servers"
	DashboardElementTypeTCBuilds       DashboardElementType = "tc-builds"
)

func (t DashboardElementType) Clean() DashboardElementType {
	switch t {
	case DashboardElementTypeReportTemplate, DashboardElementTypeGARealtime,
		DashboardElementTypeNRServers, DashboardElementTypeTCBuilds:
		return t
	default:
		return DashboardElementTypeInvalid
	}
}

func (t DashboardElementType) Validate() error {
	if t.Clean() == DashboardElementTypeInvalid {
		return ErrDashboardElementTypeInvalid
	}

	return nil
}

// Scan implements the sql.Scanner interface.
func (r *DashboardElementType) Scan(value interface{}) error {
	switch v := value.(type) {
	default:
		return fmt.Errorf("Unsupported DashboardElement type, type=%v", reflect.TypeOf(value))
	case nil:
		*r = DashboardElementTypeInvalid
	case string:
		*r = DashboardElementType(v)
	case []byte:
		*r = DashboardElementType(string(v))
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (r DashboardElementType) Value() (driver.Value, error) {
	return string(r), nil
}
