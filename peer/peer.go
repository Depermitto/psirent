package peer

import (
	"errors"
	"fmt"
	"github.com/peterh/liner"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/share"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const history = ".history"

var commands = [5]string{"handleGet", "share", "handleLs", "help", "quit"}

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
		case "handleGet":
			if len(parts) < 2 {
				fmt.Println("required positional argument <filehash> is missing")
				continue
			}

			if err = handleGet(conn, parts[1]); err != nil {
				return err
			}
		case "share":
			if len(parts) < 2 {
				fmt.Println("required positional argument <filepath> is missing")
				continue
			}

			if err = handleShare(conn, parts[1]); errors.Is(err, os.ErrNotExist) {
				fmt.Printf("inexistant file %v\n", parts[1])
			} else if errors.Is(err, share.ErrDuplicate) {
				fmt.Println("already shared")
			} else if err != nil {
				return err
			} else {
				fmt.Println(share.FileShared)
			}
		case "handleLs":
			if err = handleLs(conn); err != nil {
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
