package coordinator

import (
	"bufio"
	"errors"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/common"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"strings"
)

const storage = "coordinator/storage.json"

func CreateNetwork(addr string) error {
	// Set up the server
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalf("error creating a network on %v (%v) \n", addr, err)
	}
	defer listener.Close()
	log.Println("network created, listening for remote connections...")

	// Read from persistent storage
	vault, err := readPersistentStorage(storage)
	if err != nil {
		log.Println(err)
		return err
	}
	defer savePersistentStorage(vault, storage)
	log.Printf("storage read, %v available files...\n", len(vault))

	// Await peer connections
	conns := make(chan net.Conn)
	go func() {
		for {
			if conn, err := listener.Accept(); err == nil {
				conns <- conn
			}
		}
	}()

	// Stop the server using CTRL+D (linux/macos) or CTRL+Z (win)
	stop := make(chan struct{}, 1)
	go func() {
		_, _ = io.ReadAll(os.Stdin)
		stop <- struct{}{}
		close(stop)
	}()

	forceStop := make(chan os.Signal, 1)
	signal.Notify(forceStop, os.Interrupt)
mainloop:
	for {
		select {
		case <-forceStop:
			log.Println("stopping services...")
			break mainloop
		case <-stop:
			log.Println("gracefully stopping services...")
			break mainloop
		case conn := <-conns:
			log.Printf("peer %v connected\n", conn.RemoteAddr())
			go handlePeerConnection(vault, conn)
		}
	}
	return nil
}

func handlePeerConnection(vault map[string][]string, conn net.Conn) error {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// WARNING: Assume the client is always sending correct requests!!
		parts := strings.Split(strings.TrimSpace(scanner.Text()), ":")
		switch strings.ToLower(parts[0]) {
		case "get":
			panic("unimplemented") // TODO
		case "share":
			if err := handleShare(vault, conn, parts[1]); err == nil {
				log.Printf("peer %v shared %v\n", conn.RemoteAddr(), parts[1])
			} else if !errors.Is(err, common.ErrDuplicate) {
				return err
			}
		case "ls":
			panic("unimplemented") // TODO
		}
	}
	log.Printf("peer %v disconnected\n", conn.RemoteAddr())
	return nil
}

func handleShare(vault map[string][]string, conn net.Conn, filehash string) error {
	// We can assume we operate on IPv4
	addr := conn.RemoteAddr().String()
	if len(addr) == 0 {
		if _, err := io.WriteString(conn, common.FileNotShared); err != nil {
			return err
		}
		return common.ErrInvalidAddr
	}

	host := strings.Split(addr, ":")[0]
	if slices.Contains(vault[filehash], host) {
		if _, err := io.WriteString(conn, common.FileDuplicate); err != nil {
			return err
		}
		return common.ErrDuplicate
	}
	vault[filehash] = append(vault[filehash], host)
	_, err := io.WriteString(conn, common.FileShared)
	return err
}
