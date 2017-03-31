package types

import (
	"strconv"
)

//
// NullBool
//

type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

func NewNullBoolIfTrue(v bool) NullBool {
	if v {
		return NewNullBool(v)
	}
	return NullBool{}
}

func NewNullBool(v bool) NullBool {
	n := NullBool{}
	n.Set(v)
	return n
}

func (b *NullBool) Set(v bool) {
	b.Bool = v
	b.Valid = true
}

func (b *NullBool) Clear() {
	b.Bool = false
	b.Valid = false
}

func (b NullBool) String() string {
	if b.Valid {
		if b.Bool {
			return "true"
		} else {
			return "false"
		}
	} else {
		return "null"
	}
}

//
// NullInt32
//

type NullInt32 struct {
	Int32 int32
	Valid bool // Valid is true if Int32 is not NULL
}

func (n *NullInt32) scanString(s string) error {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return err
	}

	n.Int32 = int32(i)
	n.Valid = true
	return nil
}

func NewNullInt32(v int32) NullInt32 {
	return NullInt32{v, true}
}

func (n NullInt32) String() string {
	if n.Valid {
		return strconv.FormatInt(int64(n.Int32), 10)
	} else {
		return "null"
	}
}

//
// NullInt64
//

type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Int32 is not NULL
}

func NewNullInt64(v int64) NullInt64 {
	return NullInt64{v, true}
}

func NewNullInt64IfNotZero(v int64) NullInt64 {
	if v == 0 {
		return NullInt64{}
	} else {
		return NewNullInt64(v)
	}
}

func (n *NullInt64) Set(v int64) {
	n.Int64 = v
	n.Valid = true
}

func (n *NullInt64) Clear() {
	n.Int64 = 0
	n.Valid = false
}

func (n NullInt64) String() string {
	if n.Valid {
		return strconv.FormatInt(n.Int64, 10)
	} else {
		return "null"
	}
}

func ParseNullInt64(v string) NullInt64 {
	i, _ := strconv.ParseInt(v, 10, 64)
	if i == 0 {
		return NullInt64{}
	}
	return NewNullInt64(i)
}

//
// NullString
//

type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

func NewNullString(v string) NullString {
	n := NullString{}
	n.Set(v)
	return n
}

func (s *NullString) Set(v string) {
	s.String = v
	s.Valid = true
}

func (s *NullString) Clear() {
	s.String = ""
	s.Valid = false
}
