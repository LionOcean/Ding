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

// ServerIPInfo contains current P2P server local ips/port
type ServerIPInfo struct {
	IPv4s []string // local ipv4 list
	Port  int      `json:"port"` // unique P2P server port, depends on process id
}

func init() {
	ips, _ := localIPv4s()
	servInfo.IPv4s = ips
	servInfo.Port = os.Getpid()
	serv.Addr = fmt.Sprintf("0.0.0.0:%v", servInfo.Port)
	files_list = make([]TransferFile, 0)
}

// 作测试用，后续会删除，返回当前发送端已存的文件列表
func LogTransferFile() []TransferFile {
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

// startP2PServer start a server to received peer.
//
// Return send peer server ips and port
func StartP2PServer() (ServerIPInfo, error) {
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/download", handleDownload)
	err := serv.ListenAndServe()
	if err != nil {
		return servInfo, err
	}
	return servInfo, nil
}

// closeP2PServer close the running P2PServer
func CloseP2PServer() error {
	err := serv.Close()
	return err
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
