package transfer

import (
	"io"
	"os"
)

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
