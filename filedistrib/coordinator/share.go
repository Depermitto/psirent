package coordinator

import (
	"fmt"
	"github.com/Depermitto/psirent/filedistrib/coms"
	"github.com/Depermitto/psirent/filedistrib/persistent"
	errors2 "github.com/Depermitto/psirent/internal/errors"
	"io"
	"slices"
)

func Share(pw io.Writer, storage persistent.Storage, filehash string, peerListenAddr string) error {
	if len(peerListenAddr) == 0 {
		_, _ = fmt.Fprintln(pw, coms.ShareNotOk)
		return errors2.ErrInvalidAddr
	}
	if slices.Contains(storage[filehash], peerListenAddr) {
		_, _ = fmt.Fprintln(pw, coms.ShareDuplicate)
		return errors2.ErrShareDuplicate
	}
	storage[filehash] = append(storage[filehash], peerListenAddr)
	_, err := fmt.Fprintln(pw, coms.ShareOk)
	return err
}
