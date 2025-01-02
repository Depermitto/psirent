package peer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/peterh/liner"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/common"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const history = "peer/.history"

var commands = [5]string{"get", "share", "ls", "help", "quit"}

func Connect(addr string) error {
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		log.Fatalf("error connecting to %v (%v) \n", addr, err)
	}
	defer conn.Close()
	fmt.Printf("connected to %v\n", conn.RemoteAddr())

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

			if err = get(conn, parts[1]); err != nil {
				return err
			}
		case "share":
			if len(parts) < 2 {
				fmt.Println("required positional argument <filepath> is missing")
				continue
			}

			if err = share(conn, parts[1]); errors.Is(err, os.ErrNotExist) {
				fmt.Printf("inexistant file %v\n", parts[1])
			} else if errors.Is(err, common.ErrDuplicate) {
				fmt.Println("already shared")
			} else if err != nil {
				return err
			} else {
				fmt.Println(common.FileShared)
			}
		case "ls":
			if err = ls(conn); err != nil {
				return err
			}
		case "help":
			printHelp()
		case "quit":
			break mainloop
		default:
			fmt.Printf("unknown commad %v, please check the help command for available actions\n", cmd)
		}
	}
	return nil
}

func get(readWriter io.ReadWriter, filehash string) error {
	panic("unimplemented") // TODO
}

func ls(readWriter io.ReadWriter) error {
	panic("unimplemented") // TODO
}

func share(rw io.ReadWriter, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	filehash := sha256.Sum256(data)
	if _, err = fmt.Fprintf(rw, "SHARE:%v\n", hex.EncodeToString(filehash[:])); err != nil {
		return err
	}
	response := make([]byte, common.ResponseLength)
	if _, err = io.ReadFull(rw, response); err != nil {
		return err
	}

	if string(response) == common.FileDuplicate {
		return common.ErrDuplicate
	} else if string(response) == common.FileShared {
		return nil
	}
	return common.ErrFileNotShared
}

func printHelp() {
	fmt.Println("Commands: ")

	fmt.Println("  get <filehash>")
	fmt.Printf("    \tdownload a file from the network\n")

	fmt.Println("  share <filepath>")
	fmt.Printf("    \tdeclare a file is available for download and share them with other users\n")

	fmt.Println("  ls")
	fmt.Printf("    \tList files available for download\n")

	fmt.Println("  quit")
	fmt.Printf("    \tkill the conneciton and exit the network\n")
}
