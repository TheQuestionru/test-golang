package types

import "fmt"

var (
	ErrForbidden    = NewError("errors.Forbidden")    // "Доступ запрещен"
	ErrUnauthorized = NewError("errors.Unauthorized") // "Ошибка авторизации"
	ErrFormEmpty    = NewError("errors.FormEmpty")    // "Пустая форма"
	ErrInvalidLink  = NewError("errors.InvalidLink")  // "Неправильная ссылка"
	ErrCompanyId    = NewError("errors.CompanyId")    // "Компания заполнена неверно"
	ErrLangInvalid  = NewError("errors.LangInvalid")  // "Неверно заполнен язык"
	ErrPageNotFound = NewError("errors.PageNotFound") // "Страница не найдена"
)

func NewError(text string, args ...interface{}) *Error {
	msg := fmt.Sprintf(text, args...)
	return &Error{msg}
}

type Error struct {
	s string
}

func (e Error) Error() string {
	return e.s
}
