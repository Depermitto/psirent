package coordinator

import (
	"fmt"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/coms"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/persistent"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
	"io"
	"net"
	"strings"
	"time"
)

func Get(pw io.Writer, storage persistent.Storage, filehash string) error {

	if !get(storage, filehash) {
		_, _ = fmt.Fprintln(pw, coms.GetNotOK)
		return errors2.ErrGetFileNotShared
	}

	var validAddresses []string

	addresses := storage[filehash]

	for i := 0; i < len(addresses); {
		address := addresses[i]

		d := net.Dialer{Timeout: 1 * time.Second} // timeout
		if conn, err := d.Dial("tcp", address); err == nil {
			if Has(conn, filehash) {
				validAddresses = append(validAddresses, address)
			} else {
				persistent.Remove(storage, filehash, i)
			}

			conn.Close()
		}
		i++
	}

	if len(validAddresses) == 0 {
		_, _ = fmt.Fprintln(pw, coms.GetNoPeer)
		return errors2.ErrGetNoPeerOnline

	}

	peerList := strings.Join(validAddresses, coms.LsSeparator)

	_, err := fmt.Fprintf(pw, "%s\n", peerList)
	return err

}

func get(storage persistent.Storage, filehash string) bool {
	if _, exists := storage[filehash]; exists {
		return true
	}
	return false
}
