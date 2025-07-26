package _struct

import (
	"fmt"
	"reflect"
	"strings"
)

func StructValid[T any](data map[string]any) bool {
	if _, err := MatchingStructFromMap[T](data); err == nil {
		return true
	}
	return false
}

// MatchingStructFromMap Заполняет структуру из набора данных. Тип дженерика - структура, в которую будут записаны данные.
// Параметры:
//
//	data - набор данных
//
// Возвращаемые значения:
//
//	T - структура, в которую записаны данные. Если ошибка, то пустая структура.
//	error - ошибка, если произошла ошибка.
func MatchingStructFromMap[T any](data map[string]any) (T, error) {
	var result T
	res, err := FillingStructFromMap(result, data)
	if err != nil {
		return result, err
	}
	return *res.(*T), nil
}

// FillingStructFromMap Заполняет структуру из набора данных.
// Параметры:
//
//	target - структура, в которую будут записаны данные.
//	data - набор данных
//
// Возвращаемые значения:
//
//	T - структура (указатель), в которую записаны данные. Если ошибка, то nil.
//	error - ошибка, если произошла ошибка
func FillingStructFromMap(target any, data map[string]any) (any, error) {
	targetVal := reflect.ValueOf(target)
	// Если передано значение (не указатель)
	if targetVal.Kind() != reflect.Ptr {
		// Создаём новый указатель на копию переданного значения
		newPtr := reflect.New(targetVal.Type())
		newPtr.Elem().Set(targetVal)
		target = newPtr.Interface()
		targetVal = reflect.ValueOf(target)
	}

	// Теперь проверяем, что это указатель на структуру
	if targetVal.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("target must be a pointer to a struct or a struct value")
	}

	if err := fillingStructFromMap(data, target); err != nil {
		return nil, err
	}

	return target, nil
}

// CastStruct Проводит объект в указанный тип.
// Тип дженерика НЕ указатель.
// На входе может быть указатель или значение.
// На выходе всегда указатель.
func CastStruct[T any](obj any) (*T, bool) {
	if o, ok := obj.(T); ok {
		return &o, true
	}
	if o, ok := obj.(*T); ok {
		return o, true
	}
	return nil, false
}

func StructToMapOrFail(obj interface{}) map[string]interface{} {
	if data, err := StructToMap(obj); err != nil {
		panic(err)
	} else {
		return data
	}
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	if obj == nil {
		return nil, fmt.Errorf("nil pointer to struct")
	}

	// Получаем Value и Type переданного объекта
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	// Обработка nil указателей
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return nil, fmt.Errorf("invalid nil pointer")
	}

	// Если передали указатель, получаем значение по указателю
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// Проверяем, что это действительно структура
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", val.Kind())
	}

	result := make(map[string]interface{})

	// Итерируемся по полям структуры
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Пропускаем неэкспортируемые поля
		if fieldType.PkgPath != "" {
			continue
		}

		// Получаем имя поля из тега json или используем имя поля
		fieldName := fieldType.Name
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			// Игнорируем поля с тегом json:"-"
			if jsonTag == "-" {
				continue
			}
			// Берем первую часть тега (до запятой, если есть опции)
			fieldName = strings.Split(jsonTag, ",")[0]
		}

		// Добавляем поле в мапу
		if field.Kind() == reflect.Struct {
			// Рекурсивно обрабатываем вложенные структуры
			nestedMap, err := StructToMap(field.Interface())
			if err != nil {
				return nil, err
			}
			result[fieldName] = nestedMap
		} else {
			result[fieldName] = field.Interface()
		}
	}

	return result, nil
}

// fillingStructFromMap рекурсивно заполняет структуру из map[string]interface{}
func fillingStructFromMap(data map[string]interface{}, target interface{}) error {
	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	targetVal = targetVal.Elem()
	targetType := targetVal.Type()

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldVal := targetVal.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		// Если структура содержит json-тег, используем его, иначе используем имя поля
		fieldKey := getPureJSONTag(field)
		if fieldKey == "" {
			fieldKey = field.Name
		}

		rawValue, exists := data[fieldKey]
		if !exists {
			return fmt.Errorf("поле \"%s\" отсутствует в данных", fieldKey)
		}

		// Если поле — структура, обрабатываем рекурсивно
		if field.Type.Kind() == reflect.Struct {
			nestedMap, ok := rawValue.(map[string]interface{})
			if !ok {
				return fmt.Errorf("поле \"%s\" должно быть объектом, но получен %T", fieldKey, rawValue)
			}

			nestedPtr := reflect.New(field.Type)
			if err := fillingStructFromMap(nestedMap, nestedPtr.Interface()); err != nil {
				return err
			}

			fieldVal.Set(nestedPtr.Elem())
			continue
		}

		// Обычные поля (не структуры)
		convertedValue, ok := convertType(rawValue, field.Type)
		if !ok {
			return fmt.Errorf(
				"неверный тип для поля \"%s\": ожидается %v, получен %T",
				fieldKey, field.Type, rawValue,
			)
		}

		fieldVal.Set(convertedValue)
	}

	return nil
}

// convertType преобразует значение в нужный тип (аналогично typeTranslation)
func convertType(val interface{}, targetType reflect.Type) (reflect.Value, bool) {
	valValue := reflect.ValueOf(val)

	// Если типы совместимы, возвращаем как есть
	if valValue.Type().AssignableTo(targetType) {
		return valValue, true
	}

	// Ручное преобразование
	switch valValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if targetType.Kind() == reflect.Int {
			return reflect.ValueOf(int(val.(int64))), true
		}
	case reflect.Float32, reflect.Float64:
		if targetType.Kind() == reflect.Int {
			return reflect.ValueOf(int(val.(float64))), true
		}
	case reflect.String:
		if targetType.Kind() == reflect.Int {
			// Можно добавить парсинг строки в int, если нужно
			return reflect.Value{}, false
		}
	default:
		return reflect.Value{}, false
	}

	return reflect.Value{}, false
}

func getPureJSONTag(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return field.Name
	}

	name, _, _ := strings.Cut(jsonTag, ",")
	return name
}
