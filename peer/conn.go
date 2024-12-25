package peer

import (
	"errors"
	"fmt"
	"github.com/peterh/liner"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func Connect(addr string) error {
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		log.Fatalf("error connecting to %v (%v) \n", addr, err)
	}

	log.Printf("connected to %v\n", conn.RemoteAddr())
	err = handleConnection(conn)
	log.Println("disconnected")

	return err
}

const history = ".history"

// handleConnection is a dummy method that echos back any data that is sent to it
func handleConnection(conn net.Conn) error {
	defer conn.Close()

	commands := [5]string{"get", "share", "ls", "help", "quit"}
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
	}

	quit := false
	for !quit {
		if cmd, err := l.Prompt("psirent> "); err == nil {
			cmd = strings.TrimSpace(cmd)
			now := time.Now().Format("02/01/2006 15:04:05\t")
			l.AppendHistory(now + cmd)
			switch cmd {
			case "get":
				// TODO get
			case "share":
				// TODO share
			case "ls":
				// TODO ls
			case "help":
				fmt.Println("Commands: ")

				fmt.Println("  get <file-id>")
				fmt.Printf("    \tdownload a file from the network\n")

				fmt.Println("  share <file-path>")
				fmt.Printf("    \tdeclare a file is available for download and share them with other users\n")

				fmt.Println("  ls")
				fmt.Printf("    \tList files available for download\n")

				fmt.Println("  quit")
				fmt.Printf("    \tkill the conneciton and exit the network\n")
			case "quit":
				fmt.Println("Goodbye 🙂")
				quit = true
			default:
				fmt.Printf("unknown commad %v\n, please check the help command for available actions", cmd)
			}
		} else if errors.Is(err, io.EOF) {
			quit = true
		} else if errors.Is(err, liner.ErrPromptAborted) {
			os.Exit(1)
		}
	}

	if f, err := os.Create(history); err == nil {
		l.WriteHistory(f)
		f.Close()
	}
	return nil
}
