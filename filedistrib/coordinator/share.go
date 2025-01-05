package coordinator

import (
	"fmt"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/coms"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/persistent"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
	"io"
	"slices"
)

func Share(pw io.Writer, storage persistent.Storage, filehash string, peerListenAddr string) error {
	// We can assume we operate on IPv4
	if len(peerListenAddr) == 0 {
		if _, err := fmt.Fprintln(pw, coms.ShareNotOk); err != nil {
			return err
		}
		return errors2.ErrInvalidAddr
	}

	if slices.Contains(storage[filehash], peerListenAddr) {
		if _, err := fmt.Fprintln(pw, coms.ShareDuplicate); err != nil {
			return err
		}
		return errors2.ErrShareDuplicate
	}
	storage[filehash] = append(storage[filehash], peerListenAddr)
	_, err := fmt.Fprintln(pw, coms.ShareOk)
	return err
}
