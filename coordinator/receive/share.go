package receive

import (
	"fmt"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/coordinator/storage"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/coms"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
	"io"
	"slices"
)

func Share(pw io.Writer, s storage.Storage, filehash string, peerListenAddr string) error {
	// We can assume we operate on IPv4
	if len(peerListenAddr) == 0 {
		if _, err := fmt.Fprintln(pw, coms.ShareNotOk); err != nil {
			return err
		}
		return errors.ErrInvalidAddr
	}

	if slices.Contains(s[filehash], peerListenAddr) {
		if _, err := fmt.Fprintln(pw, coms.ShareDuplicate); err != nil {
			return err
		}
		return errors.ErrShareDuplicate
	}
	s[filehash] = append(s[filehash], peerListenAddr)
	_, err := fmt.Fprintln(pw, coms.ShareOk)
	return err
}
