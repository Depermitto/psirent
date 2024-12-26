package coordinator

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func HandleShare(vault map[string][]string, conn net.Conn, filehash string) error {
	// We can assume we operate on IPv4
	addr := conn.RemoteAddr().String()
	if len(addr) == 0 {
		return fmt.Errorf("invalid addr %v", addr)
	}

	host := strings.Split(addr, ":")[0]
	vault[filehash] = append(vault[filehash], host)

	// NOTE: we could send EOF here
	_, err := io.WriteString(conn, "OK")
	return err
}
