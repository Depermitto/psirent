package peer

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/coms"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/persistent"
	errors2 "gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/errors"
)

func Get(crw io.ReadWriter, filehash string, storage persistent.Storage) (err error) {
	// Send
	if _, err = fmt.Fprintf(crw, "GET:%v\n", filehash); err != nil {
		return
	}
	// Wait
	scanner := bufio.NewScanner(crw)
	if !scanner.Scan() {
		err = errors2.ErrLostConnection
		return
	}

	response := scanner.Text()
	if response == coms.GetNoPeer {
		err = errors2.ErrGetNoPeerOnline
	} else if response == coms.GetNotOK {
		err = errors2.ErrGetFileNotShared
	} else {
		filename := filehash
		peerList := strings.Split(response, coms.LsSeparator)
		fmt.Printf("Peers: %v\n", peerList) // Temporary log
		// Get the fragments
		total_fragments := len(peerList)
		for i := 1; i <= total_fragments; i++ {
			fmt.Println("Peer", peerList[i-1])
			conn, err := net.Dial("tcp", peerList[i-1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			// Send
			if _, err = fmt.Fprintf(conn, "FRAG:%v:%v:%v\n", i, total_fragments, filehash); err != nil {
				return err
			}
			// Wait
			reader := bufio.NewReader(conn)
			// Read filename
			filename, err = reader.ReadString('\n')
			if err != nil {
				return err
			}
			filename = strings.TrimSuffix(filename, "\n")
			fmt.Println("Filename:", filename)

			// Try to read status
			status, err := reader.ReadString('\n')
			if err == nil {
				// got status
				msg := strings.TrimSuffix(status, "\n")
				if msg == coms.FragEmpty {
					fmt.Println("Fragment empty")
					continue
				} else if msg == coms.FragNotOk {
					fmt.Println("Fragment not ok")
					continue
				}
			} else if err != io.EOF {
				fmt.Println(err)
				return err
			}
			// assume that status is just file content, will be written to file
			// Save fragment
			file, err := os.Create(fmt.Sprintf("%s.frag%d", filehash, i))
			if err != nil {
				return err
			}
			fmt.Println("Saving fragment...")
			file.Write([]byte(status))
			buffer := make([]byte, 4096) // Adjust buffer size as needed
			for {
				bytesRead, err := reader.Read(buffer)
				if err != nil {
					if err == io.EOF {
						break // End of file, stop reading
					}
					return err // Handle other errors
				}
				if bytesRead > 0 {
					if _, writeErr := file.Write(buffer[:bytesRead]); writeErr != nil {
						return writeErr
					}
				}
			}
			fmt.Printf("Fragment %d saved\n", i)
			conn.Close()
			file.Close()
		}
		// Check if the result folder exists and create it if not
		if _, err = os.Stat("received"); os.IsNotExist(err) {
			err = os.Mkdir("received", 0755)
			if err != nil {
				return err
			}
		}
		fmt.Printf("Reassembling the file %v...\n", filename)
		// Reassemble the file
		filename = "received/" + filename
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer file.Close()
		for i := 1; i <= total_fragments; i++ {
			fileFragment, err := os.Open(fmt.Sprintf("%s.frag%d", filehash, i))
			if err != nil {
				return err
			}
			// Append fragment to file
			_, err = io.Copy(file, fileFragment)
			if err != nil {
				return err
			}
			fileFragment.Close()
			// Remove fragment
			err = os.Remove(fmt.Sprintf("%s.frag%d", filehash, i))
			if err != nil {
				return err
			}
		}
		fmt.Println("File reassembled!")
		fmt.Println("Sharing the file...")
		// Share the file
		err = HandleShare(crw, filename, storage)
		if err != nil {
			return err
		}

	}
	return
}
