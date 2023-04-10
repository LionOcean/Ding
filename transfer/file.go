package transfer

import (
	"io"
	"os"
)

// use 40M temperary cache to store instead of bufio.NewReader, NewReader().Read could only read one byte every time, it's too slow although memory friendly.
const MAX_BUFFER_BYTE = 1024 * 1024 * 40

// readByte read bytes from os.File
// if file.Stat or io.ReadAt failed, return an no-nil error
// if file size is greater than MAX_BUFFER_BYTE, start byte piece transfer
func readByte(file *os.File, offset int64) ([]byte, bool, error) {
	fileInfo, err := file.Stat()
	size := fileInfo.Size()
	isChunked := false
	var buf []byte
	// only too large bytes file needs to be splited
	if size > MAX_BUFFER_BYTE {
		buf = make([]byte, MAX_BUFFER_BYTE)
		isChunked = true
	} else {
		buf = make([]byte, size)
		offset = 0
	}
	if err != nil {
		return []byte{}, isChunked, err
	}
	n, err := file.ReadAt(buf, offset)
	// reached the file end bytes, you need reurn buf although buf length is greater than n
	if err == io.EOF {
		return buf[:n], isChunked, nil
	}
	if err != nil {
		return []byte{}, isChunked, err
	}
	// you should awalys get buf[:n] splice, for in one case that if the n is less than len(buf), the buf would have lots of empty byte
	return buf[:n], isChunked, nil
}

// appendFile return os.File created by append mode
// if os.OpenFile failed, return an no-nil error
func appendFile(path string) (*os.File, error) {
	// file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// you could only use create with ReadAt method, since openFile catch error
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// bufPieces depend on base, spread total to [start, end] offset splice
func bufPieces(total, base int) [][]int {
	s := make([][]int, 0)
	for i := 0; i < total; i += base {
		s = append(s, []int{i, i + base})
	}
	return s
}
