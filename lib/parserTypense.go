package helper

import "time"

func GetString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func GetInt(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return 0
}

func GetInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key].(float64); ok {
		return int64(v)
	}
	return 0
}

func GetBool(m map[string]interface{}, key string) bool {
	switch v := m[key].(type) {
	case bool:
		return v
	case string:
		return v == "true"
	}
	return false
}

func GetTimePtr(m map[string]interface{}, key string) *time.Time {
	if v, ok := m[key].(float64); ok && v > 0 {
		t := time.Unix(int64(v), 0)
		return &t
	}
	return nil
}
