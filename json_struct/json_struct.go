package json_struct

import (
	"encoding/json"
	"fmt"
	"github.com/arhitov/goutils/struct"
	"reflect"
)

func Valid(jsonStr string) bool {
	return json.Valid([]byte(jsonStr))
}

// FindMatchingStruct принимает JSON-строку и список структур (в виде пустых экземпляров).
// Возвращает первую структуру, в которую удалось декодировать данные, или ошибку, если ни одна не подошла.
// Параметры:
//
//	jsonStr - Строка json
//	structs - Структуры (не указатель)
//
// Возвращаемые значения:
//
//	interface{} - Заполненная структура (значение) или nil
//	error - ошибка если не удалось заполнить структуру
func FindMatchingStruct(jsonStr string, structs ...interface{}) (interface{}, error) {
	if !Valid(jsonStr) {
		return nil, fmt.Errorf("json not valid: %+v", jsonStr)
	}
	for _, s := range structs {
		// Проверяем, что переданный аргумент — структура (не указатель!)
		val := reflect.ValueOf(s)
		if val.Kind() == reflect.Ptr {
			return nil, fmt.Errorf("ожидается значение, получен указатель %T", s)
		} else if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("ожидается структура, получен %T", s)
		}

		// Создаём новый указатель на структуру того же типа
		newPtr := reflect.New(val.Type()).Interface()

		// Пытаемся декодировать JSON в эту структуру
		_, err := MatchingStruct(jsonStr, newPtr)
		if err == nil {
			// Возвращаем значение (разыменовываем указатель)
			return reflect.ValueOf(newPtr).Elem().Interface(), nil
		}
	}

	return nil, fmt.Errorf("ни одна из переданных структур не подошла")
}

// MatchingStruct Проверят что переданный JSON соответствует структуре.
// Параметры:
//
//	jsonStr - Строка json
//	target - Указатель на структуру или структура, которая используются для сопоставления и заполнения значениями из jsonStr
//
// Возвращаемые значения:
//
//	interface{} - Заполненная структура (указатель) или nil
//	error - ошибка если не удалось заполнить структуру
func MatchingStruct(jsonStr string, target interface{}) (interface{}, error) {
	if !Valid(jsonStr) {
		return nil, fmt.Errorf("json not valid: %+v", jsonStr)
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &rawData); err != nil {
		return nil, err
	}

	if result, err := _struct.FillingStructFromMap(target, rawData); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

// DecodeStruct Аналогично* MatchingStruct, только в виде дженерика
// * - первое значение ответа всегда структура, даже если ошибка произошла
func DecodeStruct[T any](jsonStr string) (T, error) {
	var result T

	res, err := MatchingStruct(jsonStr, result)
	if err != nil {
		return result, err
	}
	return *res.(*T), nil
}
