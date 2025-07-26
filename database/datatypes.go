package database

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

func IsEmptyJson(value json.RawMessage) bool {
	return len(value) == 0 || string(value) == "{}" || string(value) == "null"
}

// Deref (Dereference) разыменование
func Deref[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

func Ref[T string | int | float32 | float64](value T) *T {
	switch v := any(value).(type) {
	case string:
		if v == "" {
			return nil
		}
	case int, float32, float64:
		if v == 0 {
			return nil
		}
	default:
		panic(fmt.Errorf("unsupported type: %T", value))
	}
	return &value
}

type UUIDInterface interface {
	String() string
}

func UuidIsEmpty(UUID UUIDInterface) bool {
	// Проверка на nil интерфейса
	if UUID == nil {
		return true
	}

	// Проверка на uuid.Nil (если передано значение uuid.UUID)
	if UUID == uuid.Nil {
		return true
	}

	// Специальная проверка для *uuid.UUID nil pointer
	if ptr, ok := UUID.(*uuid.UUID); ok && ptr == nil {
		return true
	}

	// Проверка нулевого UUID строкой
	return UUID.String() == "00000000-0000-0000-0000-000000000000"
}

func StringToUuid(UUID string) uuid.UUID {
	return uuid.MustParse(UUID)
}

func UuidToString(UUID *uuid.UUID) string {
	if UuidIsEmpty(UUID) {
		return ""
	} else {
		return UUID.String()
	}
}
