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
)

const history = ".history"

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

appLoop:
	for {
		cmd, err := l.Prompt("psirent> ")
		if errors.Is(err, io.EOF) {
			break appLoop
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

			if err = Get(conn, parts[1]); err != nil {
				return err
			}
		case "share":
			if len(parts) < 2 {
				fmt.Println("required positional argument <filepath> is missing")
				continue
			}

			if err = Share(conn, parts[1]); err != nil {
				return err
			}
		case "ls":
			if err = Ls(conn); err != nil {
				return err
			}
		case "help":
			PrintHelp()
		case "quit":
			break appLoop
		default:
			fmt.Printf("unknown commad %v, please check the help command for available actions\n", cmd)
		}
	}
	return nil
}
