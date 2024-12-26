package peer

import "fmt"

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
