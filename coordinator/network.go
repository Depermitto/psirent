package coordinator

import (
	"bufio"
	"errors"
	"fmt"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/coordinator/receive"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/coordinator/storage"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

const storagePath = "coordinator/storage.json"

func CreateNetwork(addr string, peerListenAddr string) error {
	// Set up the server
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalf("error creating a network on %v (%v) \n", addr, err)
	}
	defer listener.Close()
	log.Println("network created, listening for remote connections...")

	// Read from persistent storage
	s, err := storage.Read(storagePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer storage.Save(s, storagePath)
	log.Printf("storage read, %v available files...\n", len(s))

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
			go handlePeerConnection(conn, s, peerListenAddr)
		}
	}
	return nil
}

func handlePeerConnection(conn net.Conn, s storage.Storage, peerListenAddr string) error {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// WARNING: Assume the client is always sending correct requests!!
		parts := strings.Split(strings.TrimSpace(scanner.Text()), ":")
		switch strings.ToLower(parts[0]) {
		case "get":
			panic("unimplemented") // TODO
		case "share":
			if err := receive.Share(conn, s, parts[1], peerListenAddr); err == nil {
				log.Printf("peer %v shared %v\n", conn.RemoteAddr(), parts[1])
			} else if !errors.Is(err, errors2.ErrShareDuplicate) {
				return err
			}
		case "ls":
			if available, err := receive.Ls(conn, s); err == nil {
				fmt.Printf("%v available files\n", available)
			}
		}
	}
	log.Printf("peer %v disconnected\n", conn.RemoteAddr())
	return nil
}
