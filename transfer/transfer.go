package transfer

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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
		err := fmt.Errorf((err.Error()))
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
	fmt.Printf("%v", resp.Body)
	err_2 := writeFileByStep(localPath, resp.Body)
	if err_2 != nil {
		return err_2
	}
	defer resp.Body.Close()
	return nil
}
