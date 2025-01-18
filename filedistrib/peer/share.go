package peer

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/coms"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/persistent"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
)

func HandleShare(conn io.ReadWriter, filepath string, storage persistent.Storage) (err error) {
	filehash, err := Share(conn, filepath)
	if os.IsNotExist(err) {
		fmt.Printf("file %v does not exist\n", filepath)
	} else if errors.Is(err, errors2.ErrShareDuplicate) {
		fmt.Println("already shared")
	} else if _, isPathErr := err.(*os.PathError); isPathErr {
		fmt.Println("can only share files, directories are not supported")
	} else if err != nil {
		return err
	} else {
		storage[filehash] = append(storage[filehash], filepath)
	}
	return nil
}

func Share(crw io.ReadWriter, filepath string) (filehash string, err error) {
	fmt.Println("Sharing", filepath)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	filehash = hex.EncodeToString(sum[:])
	// Send
	if _, err = fmt.Fprintf(crw, "SHARE:%v\n", filehash); err != nil {
		return
	}
	// Wait
	scanner := bufio.NewScanner(crw)
	if !scanner.Scan() {
		err = errors2.ErrLostConnection
		return
	}

	response := scanner.Text()
	if response == coms.ShareDuplicate {
		err = errors2.ErrShareDuplicate
	} else if response == coms.ShareNotOk {
		err = errors2.ErrShareFileNotShared
	} else if response == coms.ShareOk {
		err = nil
	} else {
		err = errors2.ErrUnknownResponse
	}
	return
}
