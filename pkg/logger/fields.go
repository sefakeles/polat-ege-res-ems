package logger

import (
	"fmt"
	"time"
)

// field implements the Field interface
type field struct {
	key   string
	value any
}

func (f field) Key() string {
	return f.key
}

func (f field) Value() any {
	return f.value
}

// Field constructors
func String(key, value string) Field {
	return field{key: key, value: value}
}

func Int(key string, value int) Field {
	return field{key: key, value: value}
}

func Int8(key string, value int8) Field {
	return field{key: key, value: value}
}

func Int16(key string, value int16) Field {
	return field{key: key, value: value}
}

func Int32(key string, value int32) Field {
	return field{key: key, value: value}
}

func Int64(key string, value int64) Field {
	return field{key: key, value: value}
}

func Uint(key string, value uint) Field {
	return field{key: key, value: value}
}

func Uint8(key string, value uint8) Field {
	return field{key: key, value: value}
}

func Uint16(key string, value uint16) Field {
	return field{key: key, value: value}
}

func Uint32(key string, value uint32) Field {
	return field{key: key, value: value}
}

func Uint64(key string, value uint64) Field {
	return field{key: key, value: value}
}

func Float32(key string, value float32) Field {
	return field{key: key, value: value}
}

func Float64(key string, value float64) Field {
	return field{key: key, value: value}
}

func Bool(key string, value bool) Field {
	return field{key: key, value: value}
}

func Time(key string, value time.Time) Field {
	return field{key: key, value: value}
}

func Duration(key string, value time.Duration) Field {
	return field{key: key, value: value}
}

func Err(err error) Field {
	if err == nil {
		return field{key: "error", value: nil}
	}
	return field{key: "error", value: err.Error()}
}

func Any(key string, value any) Field {
	return field{key: key, value: value}
}

func Stringer(key string, value fmt.Stringer) Field {
	if value == nil {
		return field{key: key, value: nil}
	}
	return field{key: key, value: value.String()}
}
