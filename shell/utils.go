package shell

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Yulian302/qugopy/models"
)

func splitCommandLine(input string) []string {
	var args []string
	var buf strings.Builder
	var inQuote rune
	escaped := false

	for _, r := range input {
		switch {
		case escaped:
			buf.WriteRune(r)
			escaped = false
		case r == '\\':
			escaped = true
		case inQuote != 0:
			if r == inQuote {
				inQuote = 0
			} else {
				buf.WriteRune(r)
			}
		case r == '"' || r == '\'':
			inQuote = r
		case r == ' ' || r == '\t':
			if buf.Len() > 0 {
				args = append(args, buf.String())
				buf.Reset()
			}
		default:
			buf.WriteRune(r)
		}
	}

	if buf.Len() > 0 {
		args = append(args, buf.String())
	}
	return args
}

func parseArgs(line string) map[string]string {
	args := map[string]string{}
	tokens := splitCommandLine(line)

	for i := 0; i < len(tokens); i++ {
		if strings.HasPrefix(tokens[i], "--") {
			key := strings.TrimPrefix(tokens[i], "--")
			if i+1 < len(tokens) && !strings.HasPrefix(tokens[i+1], "--") {
				args[key] = tokens[i+1]
				i++
			} else {
				args[key] = "" // handle flag without value
			}
		}
	}
	return args
}

func parseTaskFromCmd(line string) (models.Task, error) {
	args := parseArgs(line)

	task := models.Task{}
	taskValue := reflect.ValueOf(&task).Elem()
	taskType := taskValue.Type()

	for i := 0; i < taskValue.NumField(); i++ {
		field := taskType.Field(i)
		jsonTag := field.Tag.Get("json")
		if !field.IsExported() {
			continue
		}

		rawValue, ok := args[jsonTag]
		if !ok {
			continue
		}

		fieldValue := taskValue.Field(i)
		fieldType := field.Type

		switch fieldType.Kind() {
		case reflect.String:
			fieldValue.SetString(rawValue)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.ParseInt(rawValue, 10, 64); err == nil {
				fieldValue.SetInt(intVal)
			} else {
				return models.Task{}, fmt.Errorf("⚠️ Int parse error for %s: %v\n", jsonTag, err)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if uintVal, err := strconv.ParseUint(rawValue, 10, 64); err == nil {
				fieldValue.SetUint(uintVal)
			} else {
				return models.Task{}, fmt.Errorf("⚠️ Uint parse error for %s: %v\n", jsonTag, err)
			}

		default:
			if fieldType == reflect.TypeOf(json.RawMessage{}) {
				fieldValue.Set(reflect.ValueOf(json.RawMessage(rawValue)))
			} else {
				return models.Task{}, fmt.Errorf("⚠️ Unknown field type for %s: %v\n", jsonTag, fieldType)
			}
		}
	}

	return task, nil
}
