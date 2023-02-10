package transfer

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	serv       http.Server // P2P server
	servInfo   ServerIPInfo
	files_list []TransferFile // file list that would be received
)

// TransferFile is the file that send to received peer
type TransferFile struct {
	Path string `json:"path"` // file absolute path
	Name string `json:"name"` // file base name with dot
	Size int    `json:"size"` // file byte length
}

// ServerIPInfo contains current P2P server local ips/port
type ServerIPInfo struct {
	IPv4s []string // local ipv4 list
	Port  int      `json:"port"` // unique P2P server port, depends on process id
}

func init() {
	ips, _ := localIPv4s()
	servInfo.IPv4s = ips
	servInfo.Port = os.Getpid()
	serv.Addr = fmt.Sprintf("0.0.0.0:%d", servInfo.Port)
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
// files could be multiple.
func RemoveTransferFiles(files ...TransferFile) {
	files_list = splice(files_list, func(v TransferFile, i int) bool {
		return some(files, func(t TransferFile, all []TransferFile, j int) bool {
			return strings.Compare(v.Path, t.Path) == 0
		})
	})
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

// ServerIPAddr return send peer server ips and port
func ServerIPAddr() ServerIPInfo {
	return servInfo
}

// DownloadFile make a GET request to remotePath, next write response.body to local file with buffer pieces.
//
// remotePath must be a full http url path eg: http://192.168.0.0.1:8090/download?path=files://user/ding/app.go
//
// localPath must be a absolute path eg: files://user/ding/new_app.go
func DownloadFile(remotePath, localPath string) error {
	resp, err := http.Get(remotePath)
	if err != nil {
		return err
	}
	err_2 := writeFileByStep(localPath, resp.Body)
	if err_2 != nil {
		return err_2
	}
	defer resp.Body.Close()
	return nil
}
