package a3interface

import (
	"fmt"
	"strings"
)

func escapeForSQF(str string) string {
	return strings.ReplaceAll(str, `"`, `""`)
}

func ToArmaHashMap(data interface{}) string {
	switch v := data.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, escapeForSQF(v))
	case int, int32, int64, float32, float64, bool:
		return fmt.Sprintf(`%v`, v)
	case map[string]interface{}:
		return toArmaHashMapMapStringInterface(v)
	case []map[string]interface{}:
		return toArmaHashMapMapStringInterfaceArray(v)
	case map[string]string:
		return toArmaHashMapMapStringString(v)
	case []interface{}:
		return toArmaHashMapInterfaceArray(v)
	default:
		return fmt.Sprintf(`"%s"`, escapeForSQF(fmt.Sprintf("%v", data)))
	}
}

func toArmaHashMapInterfaceArray(data []interface{}) string {
	var items []string
	for _, item := range data {
		items = append(items, ToArmaHashMap(item))
	}
	return "[" + strings.Join(items, ", ") + "]"
}

func toArmaHashMapMapStringInterface(data map[string]interface{}) string {
	var pairs []string
	for key, value := range data {
		pairs = append(pairs, fmt.Sprintf(`["%s", %s]`, escapeForSQF(key), ToArmaHashMap(value)))
	}
	return "[" + strings.Join(pairs, ", ") + "]"
}

func toArmaHashMapMapStringInterfaceArray(data []map[string]interface{}) string {
	var maps []string
	for _, m := range data {
		maps = append(maps, toArmaHashMapMapStringInterface(m))
	}
	return "[" + strings.Join(maps, ", ") + "]"
}

func toArmaHashMapMapStringString(data map[string]string) string {
	var pairs []string
	for key, value := range data {
		pairs = append(pairs, fmt.Sprintf(`["%s", "%s"]`, escapeForSQF(key), escapeForSQF(value)))
	}
	return "[" + strings.Join(pairs, ", ") + "]"
}
