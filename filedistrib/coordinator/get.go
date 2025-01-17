package coordinator

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"strings"
	"time"

	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/coms"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/persistent"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/constants"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
)

func Get(pw io.Writer, storage persistent.Storage, filehash string) error {
	_, ok := storage[filehash]
	if !ok {
		_, _ = fmt.Fprintln(pw, coms.GetNotOK)
		return errors2.ErrGetFileNotShared
	}

	var validAddresses []string
	// randomize the addres order for load balancing
	addresses := make([]string, len(storage[filehash]))
	perm := rand.Perm(len(addresses))
	for i, v := range(perm) {
		addresses[v] = storage[filehash][i]
	}
	for i := 0; i < len(addresses); i++ {
		address := addresses[i]
		d := net.Dialer{Timeout: 1 * time.Second} // timeout
		if conn, err := d.Dial("tcp", address); err == nil {
			if Has(conn, filehash) {
				validAddresses = append(validAddresses, address)
				// limit the number of addresses
				if len(validAddresses) == constants.MAX_ADDR_NUM {
					break
				}
			} else {
				persistent.Remove(storage, filehash, i)
			}
			conn.Close()
		}
	}

	if len(validAddresses) == 0 {
		_, _ = fmt.Fprintln(pw, coms.GetNoPeer)
		return errors2.ErrGetNoPeerOnline
	}

	peerList := strings.Join(validAddresses, coms.LsSeparator)
	_, err := fmt.Fprintf(pw, "%s\n", peerList)
	return err
}
