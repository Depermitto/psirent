package coordinator

import (
	"bufio"
	"fmt"
	"github.com/Depermitto/psirent/filedistrib/coms"
	"io"
)

func Has(prw io.ReadWriter, filehash string) bool {
	// Send
	if _, err := fmt.Fprintf(prw, "HAS:%v\n", filehash); err != nil {
		return false
	}
	// Wait
	scanner := bufio.NewScanner(prw)
	if !scanner.Scan() {
		return false
	}
	return scanner.Text() == coms.HasOk
}
