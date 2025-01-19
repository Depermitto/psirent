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
	host := flag.String("host", "localhost", "host of the coordinator")
	port := flag.Uint("port", 6000, "port the coordinator listens on")
	peerListenHost := flag.String("host-peer", "localhost", "host of the peer")
	peerListenPort := flag.Uint("port-peer", 6001, "port the peer listens on")

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

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	flag.Parse()
	addr := fmt.Sprintf("%v:%v", *host, *port)
	peerListenAddr := fmt.Sprintf("%v:%v", *peerListenHost, *peerListenPort)

	command := flag.Arg(0)
	if command == "create-network" {
		_ = filedistrib.CreateNetwork(addr)
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
