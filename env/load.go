package env

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

// LoadEnv загружает значения из переменных окружения в структуру на основе тегов env
func LoadEnv(config any) error {
	return loadEnvValue(reflect.ValueOf(config).Elem())
}

func loadEnvValue(v reflect.Value) error {
	t := v.Type()

	// Проходим по всем полям структуры
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Если поле является структурой, рекурсивно обрабатываем её
		if field.Type.Kind() == reflect.Struct {
			if err := loadEnvValue(value); err != nil {
				return err
			}
			continue
		}

		// Получаем тег env
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		// Разбираем тег на имя переменной и значение по умолчанию
		parts := strings.Split(envTag, ";")
		envVar := strings.TrimSpace(parts[0])
		defaultValue := ""
		if len(parts) > 1 {
			defaultPart := strings.TrimSpace(parts[1])
			if strings.HasPrefix(defaultPart, "default:") {
				defaultValue = strings.TrimSpace(strings.TrimPrefix(defaultPart, "default:"))
			}
		}

		// Получаем значение из переменной окружения или используем значение по умолчанию
		envValue := os.Getenv(envVar)
		if envValue == "" {
			envValue = defaultValue
		}

		// Устанавливаем значение в поле структуры в зависимости от типа
		switch value.Kind() {
		case reflect.String:
			value.SetString(envValue)
		case reflect.Int, reflect.Int64:
			if intVal, err := strconv.Atoi(envValue); err == nil {
				value.SetInt(int64(intVal))
			}
		case reflect.Bool:
			if boolVal, err := strconv.ParseBool(envValue); err == nil {
				value.SetBool(boolVal)
			}
		}
	}

	return nil
}
