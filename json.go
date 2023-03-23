package magic

import (
	"encoding/json"
	"strings"
)

func dataToMapAny(data any) map[string]any {
	b := dataToJSONBytes(data)
	m := map[string]any{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}
	return m
}

func dataToJSONBytes(data any) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return b
}

func jsonGetPath(data map[string]any, path string) any {
	if strings.Contains(path, ".") {
		return jsonGetPathRecursive(data, strings.Split(path, "."))
	}
	return data[path]
}

func jsonGetPathRecursive(data map[string]any, path []string) any {
	if len(path) == 0 {
		return nil
	}
	v := data[path[0]]
	switch m := v.(type) {
	case map[string]any:
		return jsonGetPathRecursive(m, path[1:])
	default:
		return v
	}
}
