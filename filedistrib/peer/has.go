package peer

import (
	"fmt"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/coms"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/persistent"
	"io"
	"os"
)

func Has(cw io.Writer, storage persistent.Storage, filehash string) {
	filepaths, ok := storage[filehash]
	if ok && len(filepaths) > 0 {
		for i, filepath := range filepaths {
			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				storage[filehash][i] = storage[filehash][len(storage[filehash])-1]
				storage[filehash] = storage[filehash][:len(storage[filehash])-1]

				if len(storage[filehash]) == 0 {
					delete(storage, filehash)
				}
			} else if err == nil {
				fmt.Fprintln(cw, coms.HasOk)
				return
			}
		}
	}
	fmt.Fprintln(cw, coms.HasNotOk)
}
