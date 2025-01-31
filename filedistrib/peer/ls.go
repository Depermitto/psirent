package peer

import (
	"bufio"
	"fmt"
	"github.com/Depermitto/psirent/filedistrib/coms"
	errors2 "github.com/Depermitto/psirent/internal/errors"
	"io"
	"strings"
)

func Ls(crw io.ReadWriter) ([]string, error) {
	// Send
	if _, err := fmt.Fprintln(crw, "LS"); err != nil {
		return nil, err
	}

	// Wait
	scanner := bufio.NewScanner(crw)
	if !scanner.Scan() {
		return nil, errors2.ErrLostConnection
	}
	response := scanner.Text()
	if len(response) == 0 {
		return nil, errors2.ErrLsEmpty
	}

	filehashes := strings.Split(response, coms.LsSeparator)
	return filehashes, nil
}
