package peer

import (
	"fmt"
	"io"
	"os"

	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/coms"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/filedistrib/persistent"
	"gitlab-stud.elka.pw.edu.pl/psi54/psirent/internal/constants"
)

func Fragment(cw io.Writer, storage persistent.Storage, fragNo int64, totalFragments int64, filehash string) {
	filepaths, ok := storage[filehash]
	if !ok {
		fmt.Fprintln(cw, coms.FragNotOk)
		return
	}

	filepath := filepaths[0]
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		fmt.Fprintln(cw, coms.FragNotOk)
		return
	}

	// calculate fragment size
	fileSize := fileInfo.Size()
	fragSize := fileSize / totalFragments
	if remainder := fileSize % constants.FILE_CHUNK; remainder > 0 {
		fragSize += constants.FILE_CHUNK - remainder // round up to multiple of FILE_CHUNK
	}
	fragStart := (fragNo - 1) * fragSize // fragment starting point

	// write the filename
	fmt.Fprintln(cw, filepath)
	// if fragment is determined as empty, return that information
	if fragStart > fileSize {
		fmt.Fprintln(cw, coms.FragEmpty)
		return
	}

	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintln(cw, coms.FragNotOk)
		return
	}
	defer file.Close()

	// Seek to fragment start
	_, err = file.Seek(fragStart, io.SeekStart)
	if err != nil {
		fmt.Fprintln(cw, coms.FragNotOk)
		return
	}

	// Send the fragment in chunks
	buffer := make([]byte, constants.FILE_CHUNK)
	remaining := fragSize
	for remaining > 0 {
		chunkSize := int64(len(buffer))
		if remaining < chunkSize {
			chunkSize = remaining
		}
		bytesRead, err := file.Read(buffer[:chunkSize])
		if err != nil && err != io.EOF {
			fmt.Fprintln(cw, coms.FragNotOk)
			return
		}
		if bytesRead == 0 {
			break
		}
		_, err = cw.Write(buffer[:bytesRead])
		if err != nil {
			return
		}
		remaining -= int64(bytesRead)
	}

	if closer, ok := cw.(io.Closer); ok {
		closer.Close()
	}
}
