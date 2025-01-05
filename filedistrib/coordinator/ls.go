package coordinator

import (
	"bufio"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/coms"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/persistent"
	"io"
	"net"
)

func Ls(pw io.Writer, storage persistent.Storage) (int, error) {
	bufpw := bufio.NewWriter(pw)
	available := 0
	for filehash, ips := range storage {
		for i, ip := range ips {
			if conn, err := net.Dial("tcp4", ip); err == nil {
				if Has(conn, filehash) {
					if available > 0 {
						_, _ = bufpw.WriteString(coms.LsSeparator)
					}
					_, _ = bufpw.WriteString(filehash)
					available += 1
				} else {
					storage[filehash][i] = storage[filehash][len(storage[filehash])-1]
					storage[filehash] = storage[filehash][:len(storage[filehash])-1]

					if len(storage[filehash]) == 0 {
						delete(storage, filehash)
					}
				}
				conn.Close()
			}
		}
	}
	_, _ = bufpw.WriteString("\n")
	return available, bufpw.Flush()
}
