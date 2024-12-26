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
	if _, err = fmt.Fprintf(rw, "SHARE:%v\n", hex.EncodeToString(filehash[:])); err != nil {
		return err
	}
	response := make([]byte, 2)
	if _, err = io.ReadFull(rw, response); err != nil {
		return err
	}
	if string(response) != "OK" {
		return errors.New("could not share the file")
	}
	return nil
}
