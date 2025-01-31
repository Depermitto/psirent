package peer

import (
	"fmt"
	"io"
	"os"

	"github.com/Depermitto/psirent/filedistrib/coms"
	"github.com/Depermitto/psirent/filedistrib/persistent"
	"github.com/Depermitto/psirent/internal/constants"
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
	if remainder := fragSize % constants.FileChunk; remainder > 0 {
		fragSize += constants.FileChunk - remainder // round up to multiple of FILE_CHUNK
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

	// Send the fragment size
	fmt.Fprintln(cw, min(fragSize, fileSize-fragStart))

	// Send the fragment in chunks
	buffer := make([]byte, constants.FileChunk)
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
