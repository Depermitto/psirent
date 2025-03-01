package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Depermitto/psirent/filedistrib"
	"github.com/Depermitto/psirent/internal/constants"
	errors2 "github.com/Depermitto/psirent/internal/errors"
	"github.com/Depermitto/psirent/internal/utils"
	"os"
	"syscall"
)

func main() {
	host := flag.String("host-coordinator", "localhost", "host of the coordinator")
	peerListenHost := flag.String("host-peer", "localhost", "host of the peer")

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
	addr := fmt.Sprintf("%v:%v", *host, constants.CoordinatorPort)
	peerListenAddr := fmt.Sprintf("%v:%v", *peerListenHost, constants.PeerPort)

	command := flag.Arg(0)
	if command == "create-network" {
		_ = filedistrib.CreateNetwork(addr)
	} else if command == "connect" {
	retry:
		err := filedistrib.Connect(addr, peerListenAddr)
		if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, errors2.ErrLostConnection) {
			err = utils.Reconnect(addr, peerListenAddr)
			if err != nil {
				fmt.Println("host disconnected, closing connection...")
			} else {
				goto retry // Retry Connect after a successful Reconnect
			}
		} else if err != nil {
			fmt.Println("unknown error occurred, closing connection...")
		}
	} else {
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
