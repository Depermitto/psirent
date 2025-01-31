package peer

import (
	"fmt"
	"github.com/Depermitto/psirent/filedistrib/coms"
	"github.com/Depermitto/psirent/filedistrib/persistent"
	"io"
	"os"
)

func Has(cw io.Writer, storage persistent.Storage, filehash string) {
	filepaths, ok := storage[filehash]
	if ok && len(filepaths) > 0 {
		for i, filepath := range filepaths {
			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				persistent.Remove(storage, filehash, i)
			} else if err == nil {
				fmt.Fprintln(cw, coms.HasOk)
				return
			}
		}
	}
	fmt.Fprintln(cw, coms.HasNotOk)
}
