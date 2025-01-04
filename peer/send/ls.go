package send

import (
	"bufio"
	"fmt"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/coms"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
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
