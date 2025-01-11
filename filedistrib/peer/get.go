package peer

import (
	"bufio"
	"fmt"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/coms"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
	"io"
	"strings"
)

func Get(crw io.ReadWriter, filehash string) (err error) {

	// Send
	if _, err = fmt.Fprintf(crw, "GET:%v\n", filehash); err != nil {
		return
	}
	// Wait
	scanner := bufio.NewScanner(crw)
	if !scanner.Scan() {
		err = errors2.ErrLostConnection
		return
	}

	response := scanner.Text()

	if response == coms.GetNoPeer {
		err = errors2.ErrGetNoPeerOnline
	} else if response == coms.GetNotOK {
		err = errors2.ErrGetFileNotShared
	} else {
		peerList := strings.Split(response, coms.LsSeparator)
		fmt.Printf("Peers: %v\n", peerList) // // Temporary log
		// TODO: Implement file download functionality
	}
	return
}
