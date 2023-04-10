package transfer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	EVENT_DOWN_PROGRESS = "EVENT_DOWN_PROGRESS" // donwload:progress,wails emitting events by go bingdings
)

type ReceivePeer struct {
	Ctx context.Context // wails ctx
}

func NewReceivePeer() *ReceivePeer {
	return new(ReceivePeer)
}

// ReceivingFiles return send peer uploaded files list.
//
// remoteAddr must be a ip:port format eg: 192.168.0.0.1:8090.
func (rp *ReceivePeer) ReceivingFiles(remoteAddr string) (string, error) {
	remotePath := new(strings.Builder)
	remotePath.WriteString("http://")
	remotePath.WriteString(remoteAddr)
	remotePath.WriteString("/list")
	resp, err := http.Get(remotePath.String())
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func downloadFile(url string, rangeHead []int, file *os.File, wg *sync.WaitGroup, writeBuf chan<- int) {
	defer wg.Done()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", rangeHead[0], rangeHead[1]))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	defer res.Body.Close()
	// must write data by exact offset
	file.WriteAt(data, int64(rangeHead[0]))
	writeBuf <- len(data)
}

// DownloadFile make a GET request to remoteAddr, next write response.body to local file with buffer pieces.
//
// remoteAddr must be a ip:port format eg: 192.168.0.0.1:8090.
//
// localPath must be a absolute path eg: files://user/ding/new_app.go.
func (rp *ReceivePeer) DownloadFile(remoteAddr string, remoteFile TransferFile, localPath string) error {
	remotePath := new(strings.Builder)
	remotePath.WriteString("http://")
	remotePath.WriteString(remoteAddr)
	remotePath.WriteString("/download?path=")
	remotePath.WriteString(url.QueryEscape(remoteFile.Path))

	file, err := appendFile(localPath)
	if err != nil {
		return err
	}
	defer file.Close()
	writeBuf := make(chan int) // store buf data length in every request go routine
	finishChunk := 0
	wg := new(sync.WaitGroup)

	pieces := bufPieces(remoteFile.Size, MAX_BUFFER_BYTE)
	for _, offset := range pieces {
		wg.Add(1)
		go downloadFile(remotePath.String(), offset, file, wg, writeBuf)
	}

	// syncrously close channel and waiting for routine taks finishing
	go func() {
		wg.Wait()
		// close channel could only be called in send routine
		close(writeBuf)
	}()

	for chunkSize := range writeBuf {
		finishChunk += chunkSize
		runtime.EventsEmit(rp.Ctx, EVENT_DOWN_PROGRESS, map[string]any{
			"name":     remoteFile.Name,
			"path":     remoteFile.Path,
			"chunk":    chunkSize,
			"finished": finishChunk,
			"total":    remoteFile.Size,
		})
	}
	return nil
}
