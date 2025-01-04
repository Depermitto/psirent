package peer

import (
	"errors"
	"fmt"
	"github.com/peterh/liner"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/peer/send"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const history = "peer/.history"

var commands = [5]string{"get", "share", "ls", "help", "quit"}

func Connect(addr string, peerListenAddr string) error {
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		log.Fatalf("error connecting to %v (%v) \n", addr, err)
	}
	defer conn.Close()
	fmt.Printf("connected to %v\n", conn.RemoteAddr())

	listener, err := net.Listen("tcp4", peerListenAddr)
	if err != nil {
		log.Fatalf("error creating a network on %v (%v) \n", peerListenAddr, err)
	}
	defer listener.Close()

	go func() {
		for {
			serverConn, _ := listener.Accept()
			time.Sleep(100 * time.Millisecond)
			serverConn.Close()
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
		if f, err := os.Open(history); err == nil {
			l.ReadHistory(f)
			f.Close()
		}
		defer func() {
			if f, err := os.Create(history); err == nil {
				l.WriteHistory(f)
				f.Close()
			}
		}()
	}

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

			panic("unimplemented") // TODO
		case "share":
			if len(parts) < 2 {
				fmt.Println("required positional argument <filepath> is missing")
				continue
			}

			if err = send.Share(conn, parts[1]); errors.Is(err, os.ErrNotExist) {
				fmt.Printf("nonexistant file %v\n", parts[1])
			} else if errors.Is(err, errors2.ErrShareDuplicate) {
				fmt.Println("already shared")
			} else if err != nil {
				return err
			}
		case "ls":
			if filehashes, err := send.Ls(conn); err == nil {
				fmt.Println(filehashes)
			}
		case "help":
			fmt.Println("Commands: ")

			fmt.Println("  get <filehash>")
			fmt.Printf("    \tdownload a file from the network\n")

			fmt.Println("  share <filepath>")
			fmt.Printf("    \tdeclare a file is available for download and share them with other users\n")

			fmt.Println("  ls")
			fmt.Printf("    \tList files available for download\n")

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
