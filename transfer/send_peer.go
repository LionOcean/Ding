package transfer

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	servPort = "65535" // unique P2P server port, use tcp maximum prot number
)

// TransferFile is the file that send to received peer
type TransferFile struct {
	Path string `json:"path"` // file absolute path
	Name string `json:"name"` // file base name with dot
	Size int    `json:"size"` // file byte length
}

type SendPeer struct {
	Ctx context.Context
	FileFlow
	PeerServer
}

func NewSendPeer() *SendPeer {
	sp := new(SendPeer)
	sp.files = make([]TransferFile, 0)
	sp.serv.Addr = net.JoinHostPort(net.IPv4zero.String(), servPort)
	return sp
}

// FileFlow contains several operations about TransferFiles
type FileFlow struct {
	files []TransferFile // file list that would be received
}

// List return current TransferFile list
func (fl *FileFlow) List() []TransferFile {
	return fl.files
}

// Append add new file to files_list, if file path is existed, omit it.
func (fl *FileFlow) Append(file TransferFile) {
	hasExisted := some(fl.files, func(v TransferFile, all []TransferFile, i int) bool {
		return strings.Compare(v.Path, file.Path) == 0
	})
	// only new file path should be appended
	if !hasExisted {
		fl.files = append(fl.files, file)
	}
}

// Remove delete files from files_list.
//
// files_list will be returned.
func (fl *FileFlow) Remove(files []TransferFile) []TransferFile {
	fl.files = splice(fl.files, func(v TransferFile, i int) bool {
		return some(files, func(t TransferFile, all []TransferFile, j int) bool {
			return strings.Compare(v.Path, t.Path) == 0
		})
	})
	return fl.List()
}

// PeerServer is the P2P server container, you should control the server start/end with it.
type PeerServer struct {
	serv http.Server // P2P server
}

// StartServer start a server to received peer.
func (sp *SendPeer) StartServer() error {
	// clear DefaultServeMux registration firstly
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		handleList(w, r, sp.files)
	})
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		handleDownload(w, r, sp.files)
	})

	if err := sp.serv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// CloseServer close the running P2PServer
func (sp *SendPeer) CloseServer() error {
	err := sp.serv.Close()
	return err
}

// UploadFiles show a system dialog to choose files and upload to server
func (sp *SendPeer) UploadFiles(dialogOptions runtime.OpenDialogOptions) ([]TransferFile, error) {
	emptyfiles := make([]TransferFile, 0)
	paths, err := runtime.OpenMultipleFilesDialog(sp.Ctx, dialogOptions)
	if err != nil {
		return emptyfiles, err
	}
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return emptyfiles, err
		}
		defer file.Close()
		fileInfo, err := file.Stat()
		if err != nil {
			return emptyfiles, err
		}
		sp.Append(TransferFile{
			Path: path,
			Name: fileInfo.Name(),
			Size: int(fileInfo.Size()),
		})
	}
	return sp.List(), nil
}
