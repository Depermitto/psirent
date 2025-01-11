package filedistrib

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/peterh/liner"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/peer"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/persistent"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	peerHistoryPath = ".history"
	sharedFilesPath = "filedistrib/peer/storage.json"
)

var commands = [5]string{"get", "share", "ls", "help", "quit"}

func Connect(addr string, peerListenAddr string) error {
	// Connect to the coordinator
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		log.Fatalf("error connecting to %v (%v) \n", addr, err)
	}
	defer conn.Close()
	fmt.Printf("connected to %v\n", conn.RemoteAddr())

	// Read from persistent storage
	storage, err := persistent.Read(sharedFilesPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer persistent.Save(storage, sharedFilesPath)
	log.Printf("storage read, %v available files...\n", len(storage))

	// Listen for messages from the coordinator or other peers
	listener, err := net.Listen("tcp4", peerListenAddr)
	if err != nil {
		log.Fatalf("error creating a network on %v (%v) \n", peerListenAddr, err)
	}
	defer listener.Close()

	// Await coordinator/peer connections
	go func() {
		for {
			if conn, err := listener.Accept(); err == nil {
				go handleIncomingConnection(conn, storage)
			}
		}
	}()

	// Setting up the user interface
	l := liner.NewLiner()
	defer l.Close()
	{
		l.SetTabCompletionStyle(liner.TabPrints)
		l.SetCtrlCAborts(true)
		l.SetCompleter(func(line string) (c []string) {
			for _, n := range commands {
				if strings.HasPrefix(n, strings.ToLower(line)) {
					c = append(c, n+" ")
				}
			}
			return
		})
		if f, err := os.Open(peerHistoryPath); err == nil {
			l.ReadHistory(f)
			f.Close()
		}
		defer func() {
			if f, err := os.Create(peerHistoryPath); err == nil {
				l.WriteHistory(f)
				f.Close()
			}
		}()
	}
	// Contains app loop
	return handleOutgoingConnection(conn, storage, l)
}

func handleOutgoingConnection(conn net.Conn, storage persistent.Storage, l *liner.State) error {
mainloop:
	for {
		cmd, err := l.Prompt("psirent> ")
		if errors.Is(err, io.EOF) {
			break mainloop
		} else if err != nil {
			return err
		}
		l.AppendHistory(cmd)

		parts := strings.Split(strings.TrimSpace(cmd), " ")
		switch parts[0] {
		case "get":
			if len(parts) < 2 {
				fmt.Println("required positional argument <filehash> is missing")
				continue
			}
			filehash := parts[1]
			err :=  peer.Get(conn, filehash)
			if errors.Is(err, errors2.ErrGetFileNotShared){
				fmt.Printf("no file found with the specified hash\n")
			} else if errors.Is(err, errors2.ErrGetNoPeerOnline) {
				fmt.Println("no peers that have the requested file are reachable right now")
			} else if err != nil {
				return err
			}
		case "share":
			if len(parts) < 2 {
				fmt.Println("required positional argument <filepath> is missing")
				continue
			}

			filepath := parts[1]
			filehash, err := peer.Share(conn, filepath)
			if os.IsNotExist(err) {
				fmt.Printf("nonexistant file %v\n", filepath)
			} else if errors.Is(err, errors2.ErrShareDuplicate) {
				fmt.Println("already shared")
			} else if err != nil {
				return err
			} else {
				storage[filehash] = append(storage[filehash], filepath)
			}
		case "ls":
			if filehashes, err := peer.Ls(conn); err == nil {
				fmt.Println(filehashes)
			}
		case "help":
			fmt.Println("Commands: ")

			fmt.Println("  get <filehash>")
			fmt.Printf("    \tdownload a file from the network\n")

			fmt.Println("  share <filepath>")
			fmt.Printf("    \tdeclare a file is available for download and share them with other users\n")

			fmt.Println("  ls")
			fmt.Printf("    \tlist files available for download\n")

			fmt.Println("  quit")
			fmt.Printf("    \tkill the conneciton and exit the network\n")
		case "quit":
			break mainloop
		default:
			fmt.Printf("unknown commad %v, please check the help command for available actions\n", cmd)
		}
	}
	return nil
}

func handleIncomingConnection(conn net.Conn, storage persistent.Storage) error {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// WARNING: Assume the other side is always sending correct requests!!
		parts := strings.Split(strings.TrimSpace(scanner.Text()), ":")
		switch strings.ToLower(parts[0]) {
		case "has":
			peer.Has(conn, storage, parts[1])
		}
	}
	return nil
}
