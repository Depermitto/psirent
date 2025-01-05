package main

import (
	"errors"
	"flag"
	"fmt"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib"
	"os"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Positional argument required, check -help for details")
		os.Exit(1)
	}

	flag.Usage = func() {
		// Tell your IDE to ignore these warnings, it is not worth checking them.
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])

		fmt.Fprintln(flag.CommandLine.Output(), "Arguments: ")
		fmt.Fprintln(flag.CommandLine.Output(), "  create-network")
		fmt.Fprintf(flag.CommandLine.Output(), "    \tact as the network coordinator\n")
		fmt.Fprintln(flag.CommandLine.Output(), "  connect")
		fmt.Fprintf(flag.CommandLine.Output(), "    \tact as a peer\n")

		fmt.Fprintln(flag.CommandLine.Output(), "Flags: ")
		flag.PrintDefaults()
	}

	host := flag.String("host", "localhost", "host to connect to/listen on")
	port := flag.Uint("port", 6000, "port on the machine")
	flag.Parse()

	addr := fmt.Sprintf("%v:%v", *host, *port)
	peerListenAddr := fmt.Sprintf("%v:%v", *host, *port+1)

	command := os.Args[1]
	if command == "create-network" {
		_ = filedistrib.CreateNetwork(addr, peerListenAddr)
	} else if command == "connect" {
		err := filedistrib.Connect(addr, peerListenAddr)
		if errors.Is(err, syscall.EPIPE) {
			fmt.Println("host disconnected, closing connection...")
		} else if err != nil {
			fmt.Println("unknown error occurred, closing connection...")
		}
	} else {
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
