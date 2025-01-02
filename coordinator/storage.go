package coordinator

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

func readPersistentStorage(path string) (vault map[string][]string, err error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return make(map[string][]string), nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &vault); err != nil {
		return nil, err
	}
	return
}

func savePersistentStorage(filemap map[string][]string, path string) error {
	bytes, err := json.MarshalIndent(filemap, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, os.ModePerm)
}
