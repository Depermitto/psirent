package send

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/coms"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
	"io"
	"os"
)

func Share(crw io.ReadWriter, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	filehash := sha256.Sum256(data)
	// Send
	if _, err = fmt.Fprintf(crw, "SHARE:%v\n", hex.EncodeToString(filehash[:])); err != nil {
		return err
	}
	// Wait
	scanner := bufio.NewScanner(crw)
	if !scanner.Scan() {
		return errors.ErrLostConnection
	}

	response := scanner.Text()
	if response == coms.ShareDuplicate {
		return errors.ErrShareDuplicate
	} else if response == coms.ShareNotOk {
		return errors.ErrShareFileNotShared
	} else if response == coms.ShareOk {
		return nil
	}
	fmt.Println(response)
	return errors.ErrUnknownResponse
}
