package coordinator

import (
	"errors"
	"io"
	"log"
	"net"
	"slices"
	"strings"
)

func handleShare(vault map[string][]string, conn net.Conn, filehash string) error {
	// We can assume we operate on IPv4
	addr := conn.RemoteAddr().String()
	if len(addr) == 0 {
		return errors.New("invalid addr")
	}

	host := strings.Split(addr, ":")[0]
	if !slices.Contains(vault[filehash], host) {
		vault[filehash] = append(vault[filehash], host)
		log.Printf("peer %v shared %v\n", conn.RemoteAddr(), filehash)
	}

	// NOTE: we could send EOF here
	_, err := io.WriteString(conn, "OK")
	return err
}
