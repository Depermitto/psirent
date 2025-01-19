package common

import (
	"fmt"
	"net"
	"time"

	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/constants"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
)

func Reconnect(addr string, peerListenAddr string) error {
	var conn net.Conn
	var err error

	for i := 0; i < constants.MAX_RETRY_ATTEMPTS; i++ {
		// Try dialing first
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			fmt.Printf("Dial attempt %d/5 failed: %v. Retrying in 2 seconds...\n", i+1, err)
			time.Sleep(constants.RETRY_DELAY)
			continue
		}

		// Successful connection
		conn.Close()
		fmt.Println("Successfully connected")
		return nil
	}

	return errors2.ErrRetryFailed
}
