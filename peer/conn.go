package peer

import (
	"bufio"
	"io"
	"log"
	"net"
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

// handleConnection is a dummy method that echos back any data that is sent to it
func handleConnection(conn net.Conn) error {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	_, err := io.WriteString(conn, "abc\n")
	if err != nil {
		return err
	}

	_, err = reader.ReadString('\n')
	if err != nil {
		return err
	}

	return nil

	// TODO create REPL-like user interface

	//l := liner.NewLiner()
	//defer l.Close()
	//{
	//	l.SetTabCompletionStyle(liner.TabPrints)
	//	l.SetCtrlCAborts(true)
	//	l.SetCompleter(func(line string) (c []string) {
	//		for _, n := range commands {
	//			if strings.HasPrefix(n, strings.ToLower(line)) {
	//				c = append(c, n+" ")
	//			}
	//		}
	//		return
	//	})
	//	if f, err := os.Open(history); err == nil {
	//		l.ReadHistory(f)
	//		f.Close()
	//	}
	//}
	//
	//quit := false
	//for !quit {
	//	if cmd, err := l.Prompt("psirent> "); err == nil {
	//		cmd = strings.TrimSpace(cmd)
	//		now := time.Now().Format("02/01/2006 15:04:05\t")
	//		l.AppendHistory(now + cmd)
	//		if cmd == "quit" {
	//			fmt.Println("Goodbye 🙂")
	//			quit = true
	//		} else {
	//			log.Printf("%v is not supported...\n", cmd)
	//		}
	//	} else if errors.Is(err, io.EOF) {
	//		quit = true
	//	} else if errors.Is(err, liner.ErrPromptAborted) {
	//		os.Exit(1)
	//	}
	//}
	//
	//if f, err := os.Create(history); err == nil {
	//	l.WriteHistory(f)
	//	f.Close()
	//}
	//return nil
}
