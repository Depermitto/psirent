package receive

import (
	"bufio"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/coordinator/storage"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/coms"
	"io"
	"net"
)

func Ls(pw io.Writer, s storage.Storage) (int, error) {
	bufw := bufio.NewWriter(pw)
	active := 0
	for filehash, ips := range s {
		for _, ip := range ips {
			if conn, err := net.Dial("tcp4", ip); err == nil {
				if active > 0 {
					_, _ = bufw.WriteString(coms.LsSeparator)
				}
				_, _ = bufw.WriteString(filehash)
				active += 1
				conn.Close()
			}
		}
	}
	_, _ = bufw.WriteString("\n")
	return active, bufw.Flush()
}
