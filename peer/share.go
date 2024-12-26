package peer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/share"
	"io"
	"os"
)

func handleShare(rw io.ReadWriter, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	filehash := sha256.Sum256(data)
	if _, err = fmt.Fprintf(rw, "SHARE:%v\n", hex.EncodeToString(filehash[:])); err != nil {
		return err
	}
	response := make([]byte, 2)
	if _, err = io.ReadFull(rw, response); err != nil {
		return err
	}

	if string(response) == share.FileDuplicate {
		return share.ErrDuplicate
	} else if string(response) == share.FileShared {
		return nil
	}
	return share.ErrFileNotShared
}
