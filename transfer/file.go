package transfer

import (
	"io"
	"os"
)

// use 50M temperary cache to store instead of bufio.NewReader, NewReader().Read could only read one byte every time, it's too slow although memory friendly.
const MAX_BUFFER_BYTE = 1024 * 1024 * 50

// readFileByStep read byte constantly from os.File created by given file path
// if os.Open or ReadByte failed, return an no-nil error
// you could receive file byte from byteHandle
func readFileByStep(path string, byteHandle func(currentBuf []byte, fileInfo os.FileInfo)) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	size := fileInfo.Size()
	var buf []byte
	// only too large bytes file needs to be splited
	if size > MAX_BUFFER_BYTE {
		buf = make([]byte, MAX_BUFFER_BYTE)
	} else {
		buf = make([]byte, size)
	}
	if err != nil {
		return err
	}
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			buf = nil
			return nil
		}
		if err != nil {
			buf = nil
			return err
		}
		// you should awalys get buf[:n] splice, for in one case that if the n is less than len(buf), the buf would have lots of empty byte
		byteHandle(buf[:n], fileInfo)
	}
}

// writeFileByStep write byte constantly from r to os.File created by given file path
// if os.OpenFile failed, return an no-nil error
func writeFileByStep(path string, r io.Reader) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.ReadFrom(r)
	return nil
}
