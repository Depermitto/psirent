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
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/constants"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
)

func HandleShare(conn io.ReadWriter, filepath string, myListenAddr string, storage persistent.Storage) (err error) {
	filehash, err := Share(conn, filepath, myListenAddr)
	if os.IsNotExist(err) {
		fmt.Printf("%s file %v does not exist\n", constants.PEER_PREFIX, filepath)
	} else if errors.Is(err, errors2.ErrShareDuplicate) {
		fmt.Printf("%s You have already shared this file\n", constants.HOST_PREFIX)
	} else if _, isPathErr := err.(*os.PathError); isPathErr {
		fmt.Printf("%s You can only share files, directories are not supported.\n", constants.HOST_PREFIX)
	} else if err != nil {
		return err
	} else {
		storage[filehash] = append(storage[filehash], filepath)
	}
	return nil
}

func Share(crw io.ReadWriter, filepath string, myListenAddr string) (filehash string, err error) {
	fmt.Println(constants.PEER_PREFIX, "Sharing", filepath, "...")
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	filehash = hex.EncodeToString(sum[:])
	// Send
	if _, err = fmt.Fprintf(crw, "SHARE:%v:%v\n", filehash, myListenAddr); err != nil {
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
