package peer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"syscall"

	// "gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/coms"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/constants"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/utils"
)

func Ls(crw io.ReadWriter) ([]string, error) {
	var response string


	err := utils.Retry(func() error {
		if crw == nil {
			return fmt.Errorf("crw is nil, cannot proceed")
		}

		// Send
		if _, err := fmt.Fprintln(crw, "LS"); err != nil {
			if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNREFUSED){ // Broken pipe
				if conn, ok := crw.(net.Conn); ok {
					// Get the remote address (host:port)
					remoteAddr := conn.RemoteAddr().String()
					// peerListenAddr := conn.LocalAddr().String()
					// reconnectErr := common.Connect(remoteAddr, peerListenAddr)
					
					// Reconnect
					tempConn, reconnectErr := net.Dial("tcp", remoteAddr)
					if reconnectErr != nil {
						return reconnectErr
					}
					crw = tempConn
					// if reconnectErr != nil {
					// 	return reconnectErr
					// }
				} else {
					return fmt.Errorf("crw is not a net.Conn, cannot retrieve address")
				}
			} else {
				return err
			}
		}

		// Wait
		scanner := bufio.NewScanner(crw)
		if !scanner.Scan() {
			return errors2.ErrLostConnection
		}

		response = scanner.Text() // Store the response
		fmt.Printf("%v", response)
		return nil
	}, constants.MAX_RETRY_ATTEMPTS, constants.RETRY_DELAY)

	if err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, errors2.ErrLsEmpty
	}

	filehashes := strings.Split(response, coms.LsSeparator)
	return filehashes, nil
}