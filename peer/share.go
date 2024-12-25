package peer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

func Share(rw io.ReadWriter, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("could not open file at %v (%v)\n", filepath, err)
	}

	filehash := sha256.Sum256(data)
	_, err = fmt.Fprintf(rw, "SHARE:%v\n", hex.EncodeToString(filehash[:]))
	response, err := io.ReadAll(rw)
	if string(response) != "OK" {
		return errors.New("could not share the file")
	}
	return err
}
