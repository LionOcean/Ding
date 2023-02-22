package transfer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	// wails emitting events by go bingdings
	EVENT_DOWN_PROGRESS = "EVENT_DOWN_PROGRESS"
)

var (
	serv       http.Server    // P2P server
	servPort   string         // unique P2P server port, use process id
	files_list []TransferFile // file list that would be received
)

// TransferFile is the file that send to received peer
type TransferFile struct {
	Path string `json:"path"` // file absolute path
	Name string `json:"name"` // file base name with dot
	Size int    `json:"size"` // file byte length
}

func init() {
	servPort = strconv.FormatInt(int64(os.Getpid()), 10)
	serv.Addr = fmt.Sprintf("0.0.0.0:%v", servPort)
	files_list = make([]TransferFile, 0)
}

// LogTransferFiles return current TransferFile list
func LogTransferFiles() []TransferFile {
	return files_list
}

// AppendTransferFile add new file to files_list, if file path is existed, omit it.
func AppendTransferFile(file TransferFile) {
	hasExisted := some(files_list, func(v TransferFile, all []TransferFile, i int) bool {
		return strings.Compare(v.Path, file.Path) == 0
	})
	// only new file path should be appended
	if !hasExisted {
		files_list = append(files_list, file)
	}
}

// RemoveTransferFiles remove files from files_list.
//
// files could be multiple. files_list will be returned.
func RemoveTransferFiles(files ...TransferFile) []TransferFile {
	files_list = splice(files_list, func(v TransferFile, i int) bool {
		return some(files, func(t TransferFile, all []TransferFile, j int) bool {
			return strings.Compare(v.Path, t.Path) == 0
		})
	})
	return LogTransferFiles()
}

// startP2PServer start a server to received peer.
func StartP2PServer() error {
	// clear DefaultServeMux registration firstly
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/download", handleDownload)

	if err := serv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// closeP2PServer close the running P2PServer
func CloseP2PServer() error {
	err := serv.Close()
	return err
}

// LocalIPAddr return current peer local ipv4 and port, if in received peer, you should omit port filed.
func LocalIPAddr() ([]string, error) {
	ip, port, err := localIPv4WithNetwork()
	if err != nil {
		return []string{}, err
	}
	port = servPort
	return []string{ip, port}, nil
}

// ReceivingFiles return send peer upaloded file list.
//
// remoteAddr must be a ip:port format eg: 192.168.0.0.1:8090.
func ReceivingFiles(remoteAddr string) (string, error) {
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
func DownloadFile(ctx context.Context, remoteAddr string, remoteFile TransferFile, localPath string) error {
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
		runtime.EventsEmit(ctx, EVENT_DOWN_PROGRESS, map[string]any{
			"name":     remoteFile.Name,
			"path":     remoteFile.Path,
			"chunk":    chunkSize,
			"finished": finishChunk,
			"total":    remoteFile.Size,
		})
	}
	return nil
}
