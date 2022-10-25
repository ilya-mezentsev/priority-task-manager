package settings

import (
	"encoding/json"
	"io/ioutil"
)

func Parse[T any](path string, foo *T) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	return json.Unmarshal(content, foo)
}
