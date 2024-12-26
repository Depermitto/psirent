package coordinator

import (
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/share"
	"io"
	"net"
	"slices"
	"strings"
)

func handleShare(vault map[string][]string, conn net.Conn, filehash string) error {
	// We can assume we operate on IPv4
	addr := conn.RemoteAddr().String()
	if len(addr) == 0 {
		if _, err := io.WriteString(conn, share.FileNotShared); err != nil {
			return err
		}
		return share.ErrInvalidAddr
	}

	host := strings.Split(addr, ":")[0]
	if slices.Contains(vault[filehash], host) {
		if _, err := io.WriteString(conn, share.FileDuplicate); err != nil {
			return err
		}
		return share.ErrDuplicate
	}
	vault[filehash] = append(vault[filehash], host)
	_, err := io.WriteString(conn, share.FileShared)
	return err
}
