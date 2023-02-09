package transfer

import (
	"bufio"
	"io"
	"os"
)

// readFileByStep read byte constantly from os.File created by given file path
// if os.Open failed, return an no-nil error
// you could receive file byte from fileChan
func readFileByStep(path string, w io.Writer, fileChan chan<- byte) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	for {
		singleByte, err := buf.ReadByte()
		if err != nil {
			close(fileChan)
			return nil
		}
		fileChan <- singleByte
	}
}
