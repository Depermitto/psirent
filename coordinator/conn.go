package coordinator

import (
	"bufio"
	"io"
	"log"
	"net"
)

func CreateNetwork(addr string) error {
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalf("error creating a network on %v (%v) \n", addr, err)
	}
	defer listener.Close()
	log.Println("network created, listening for remote connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("encountered %v\n", err.Error())
			continue
		}

		log.Printf("peer %v connected\n", conn.RemoteAddr())
		go func() {
			_ = handleConnection(conn)
			log.Printf("peer %v disconnected\n", conn.RemoteAddr())
		}()
	}
}

// handleConnection is a dummy method that echos back any data that is sent to it
func handleConnection(conn net.Conn) error {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	msg, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	if _, err = io.WriteString(conn, msg+"\n"); err != nil {
		return err
	}

	return nil

	// TODO accept incoming events and respond to them appropriately
}
