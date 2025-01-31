package filedistrib

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Depermitto/psirent/internal/constants"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/Depermitto/psirent/filedistrib/coordinator"
	"github.com/Depermitto/psirent/filedistrib/persistent"
	errors2 "github.com/Depermitto/psirent/internal/errors"
)

const storagePath = "coordinator.json"

func CreateNetwork(addr string) error {
	// Set up the server
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalf("error creating a network on %v (%v) \n", addr, err)
	}
	defer listener.Close()
	log.Printf("network created on %v, listening for remote connections...\n", addr)

	// Read from persistent storage
	storage, err := persistent.Read(storagePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer persistent.Save(storage, storagePath)
	log.Printf("storage read, %v available files...\n", len(storage))

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
			go handlePeerConnection(conn, storage)
		}
	}
	return nil
}

func handlePeerConnection(conn net.Conn, storage persistent.Storage) error {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// WARNING: Assume the other side is always sending correct requests!!
		parts := strings.Split(strings.TrimSpace(scanner.Text()), ":")
		switch strings.ToLower(parts[0]) {
		case "get":
			if err := coordinator.Get(conn, storage, parts[1]); err == nil {
				log.Printf("sent the list of peers with the file: %v to address: %v", parts[1], conn.RemoteAddr())
			} else if !errors.Is(err, errors2.ErrGetFileNotShared) && !errors.Is(err, errors2.ErrGetNoPeerOnline) {
				return err
			}
		case "share":
			filehash := parts[1]
			peerListenAddr := parts[2] + ":" + strconv.Itoa(constants.PeerPort)
			if err := coordinator.Share(conn, storage, filehash, peerListenAddr); err == nil {
				log.Printf("peer %v shared %v\n", conn.RemoteAddr(), filehash)
			} else if !errors.Is(err, errors2.ErrShareDuplicate) {
				return err
			}
		case "ls":
			if available, err := coordinator.Ls(conn, storage); err == nil {
				fmt.Printf("%v available files\n", available)
			}
		}
	}
	log.Printf("peer %v disconnected\n", conn.RemoteAddr())
	return nil
}
