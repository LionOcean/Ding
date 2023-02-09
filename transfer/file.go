package transfer

import (
	"bufio"
	"io"
	"os"
)

// readFileByStep read byte constantly from os.File created by given file path
// if os.Open or ReadByte failed, return an no-nil error
// you could receive file byte from byteHandle
func readFileByStep(path string, byteHandle func(currentBuf byte, fileInfo os.FileInfo)) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	for {
		singleByte, err := buf.ReadByte()
		if err != nil {
			return err
		}
		byteHandle(singleByte, fileInfo)
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
