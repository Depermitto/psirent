package filedistrib

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/peterh/liner"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/coms"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/peer"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/persistent"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/constants"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
)

const (
	peerHistoryPath = ".history"
	sharedFilesPath = "peer.json"
)

var commands = [5]string{"get", "share", "ls", "help", "quit"}

func Connect(addr string, myListenAddr string) error {
	// Connect to the coordinator
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		fmt.Printf("%s Error connecting to %v (%v) \n", constants.PEER_PREFIX, addr, err)
		return err
	}
	defer conn.Close()

	fmt.Printf("%s Connected to %v\n", constants.PEER_PREFIX, conn.RemoteAddr())

	// Read from persistent storage
	storage, err := persistent.Read(sharedFilesPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer persistent.Save(storage, sharedFilesPath)
	log.Printf("%s Storage read, %v files available...\n", constants.PeerPrefix, len(storage))

	// Listen for messages from the coordinator or other peers
	// @TODO: Limit the number of connections to constants.MAX_ADDR_NUM
	listener, err := net.Listen("tcp4", myListenAddr)
	if err != nil {
		log.Fatalf("%s Error creating a network on %v (%v) \n", constants.PeerPrefix, myListenAddr, err)
	}
	defer listener.Close()

	// Await coordinator/peer connections
	go func() {
		semaphore := make(chan struct{}, constants.MaxAddrNum)
		for {
			if incomingConn, err := listener.Accept(); err == nil {
				go handleIncomingConnection(incomingConn, semaphore, storage)
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
	return handleOutgoingConnection(conn, myListenAddr, storage, l)
}

func handleOutgoingConnection(conn net.Conn, myListenAddr string, storage persistent.Storage, l *liner.State) error {
mainloop:
	for {
		cmd, err := l.Prompt("PSIrent> ")
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
				fmt.Println(constants.PeerPrefix, "Required positional argument <filehash> is missing")
				continue
			}
			filehash := parts[1]
			err := peer.Get(conn, filehash, myListenAddr, storage)
			if errors.Is(err, errors2.ErrGetFileNotShared) {
				fmt.Println(constants.HostPrefix, "No file found with the specified hash")
			} else if errors.Is(err, errors2.ErrGetNoPeerOnline) {
				fmt.Println(constants.HostPrefix, "No peers that have the requested file are reachable right now")
			} else if err != nil {
				return err
			}
		case "share":
			if len(parts) < 2 {
				fmt.Println(constants.PeerPrefix, "Required positional argument <filepath> is missing")
				continue
			}

			filepath := parts[1]
			err := peer.HandleShare(conn, filepath, myListenAddr, storage)
			if err != nil {
				return err
			}
		case "ls":
			if filehashes, err := peer.Ls(conn); err == nil {
				fmt.Printf("%s %d available files:", constants.HostPrefix, len(filehashes))
				for _, filehash := range filehashes {
					fmt.Printf("\n  %v", filehash)
				}
				fmt.Println()
			} else {
				fmt.Printf("%v\n", err)
				return err
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

func handleIncomingConnection(conn net.Conn, peerSemaphore chan struct{}, storage persistent.Storage) error {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// WARNING: Assume the other side is always sending correct requests!!
		parts := strings.Split(strings.TrimSpace(scanner.Text()), ":")
		switch strings.ToLower(parts[0]) {
		case "has":
			peer.Has(conn, storage, parts[1])
		case "frag":
			// firstly convert numbers
			fragNo, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				fmt.Fprintln(conn, coms.FragNotOk)
				continue
			}
			totalFragments, err := strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				fmt.Fprintln(conn, coms.FragNotOk)
				continue
			}
			// take a slot in the semaphore
			peerSemaphore <- struct{}{}
			// download a fragment
			peer.Fragment(conn, storage, fragNo, totalFragments, parts[3])
			// release our slot
			<-peerSemaphore
		}
	}
	return nil
}
