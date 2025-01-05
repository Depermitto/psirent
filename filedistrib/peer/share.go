package peer

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/coms"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
	"io"
	"os"
)

func Share(crw io.ReadWriter, filepath string) (filehash string, err error) {
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
