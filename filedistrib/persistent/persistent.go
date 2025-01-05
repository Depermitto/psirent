package persistent

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Storage = map[string][]string

func Read(path string) (storage Storage, err error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return make(Storage), nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &storage); err != nil {
		return nil, err
	}
	return
}

func Save(storage Storage, path string) error {
	bytes, err := json.MarshalIndent(storage, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, os.ModePerm)
}
