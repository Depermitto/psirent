package peer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/peterh/liner"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var commands = [5]string{"get", "share", "ls", "help", "quit"}

const history = ".history"

func Connect(addr string) error {
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		log.Fatalf("error connecting to %v (%v) \n", addr, err)
	}
	defer conn.Close()

	fmt.Printf("connected to %v\n", conn.RemoteAddr())
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

	for {
		cmd, err := l.Prompt("psirent> ")
		if errors.Is(err, io.EOF) {
			return nil
		} else if errors.Is(err, liner.ErrPromptAborted) {
			os.Exit(1)
		} else if err != nil {
			fmt.Println("unknown error occurred")
			os.Exit(2)
		}

		parts := strings.Split(strings.TrimSpace(cmd), " ")
		l.AppendHistory(cmd)
		switch parts[0] {
		case "get":
			if len(parts) < 2 {
				fmt.Println("required positional argument <filehash> missing")
				continue
			}

			filehash := parts[1]
			if _, err = fmt.Fprintf(conn, "GET:%v\n", filehash); err != nil {
				return err
			}
		case "share":
			if len(parts) < 2 {
				fmt.Println("required positional argument <filepath> missing")
				continue
			}

			filepath := parts[1]
			data, err := os.ReadFile(filepath)
			if err != nil {
				fmt.Printf("could not open file at %v (%v)\n", filepath, err)
			}

			filehash := sha256.Sum256(data)
			if _, err = fmt.Fprintf(conn, "SHARE:%v\n", hex.EncodeToString(filehash[:])); err != nil {
				return err
			}
		case "ls":
			if _, err = fmt.Fprintln(conn, "LS"); err != nil {
				return err
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
			return nil
		default:
			fmt.Printf("unknown commad %v, please check the help command for available actions\n", cmd)
		}
	}
	return nil
}
