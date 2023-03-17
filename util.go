package magic

import (
	"encoding/base64"
	"encoding/json"
)

func Must[T any](data T, err error) T {
	if err != nil {
		panic(err)
	}
	return data
}

func dataToBase64(data any) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return []byte(base64.RawStdEncoding.EncodeToString(b)), nil
}
