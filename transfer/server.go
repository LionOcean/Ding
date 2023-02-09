package transfer

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	serv       http.Server    // P2P server
	files_list []transterFile // file list that would be received
)

const (
	responseCodeSuccess = iota + 200 // successful code
	responseCodeFail                 // faliure code
)

// transterFile is the file that send to received peer
type transterFile struct {
	Path string `json:"path"` // file absolute path
	Size int    `json:"size"` // file byte length
}

// server response JSON schema
type responseJSON struct {
	Code         int    `json:"code"`          // server action code
	Data         any    `json:"data"`          // server returned data
	ErrorMessage string `json:"error_message"` // server returned error when code is not responseCodeSuccess(200)
}

func init() {
	files_list = make([]transterFile, 0)
	serv.Addr = "0.0.0.0:8090"
}

// handleList handle /list route.
//
// it responses send file list JSON.
func handleList(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	switch method {
	case http.MethodGet:
		resInfo := responseJSON{responseCodeSuccess, files_list, ""}
		resJSON, err := json.Marshal(resInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Length", strconv.FormatInt(int64(len(resJSON)), 10))
		w.Write(resJSON)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// handleDownload handle /download route.
//
// it responses byte stream according to route query path filed.
func handleDownload(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	switch method {
	case http.MethodGet:
		path, ok := req.URL.Query()["path"]
		filePath := path[0]
		isPathExisted := some(files_list, func(v transterFile, all []transterFile, i int) bool {
			return strings.Compare(v.Path, filePath) == 0
		})
		// path filed exist and path exist in file list
		if ok && isPathExisted {
			err := readFileByStep(filePath, func(current byte, fileInfo os.FileInfo) {
				size := fileInfo.Size()
				w.Header().Add("Content-Length", strconv.FormatInt(size, 10))
				w.Write([]byte{current})
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.Header().Add("Content-Type", "application/json")
			resInfo := responseJSON{
				responseCodeFail,
				map[string]string{},
				"",
			}
			// path filed lose
			if !ok {
				resInfo.ErrorMessage = "query path filed is necessary."
			}
			// path not exist in file list
			if !isPathExisted {
				resInfo.ErrorMessage = "query path has been deleted by send peer."
			}
			resJSON, err := json.Marshal(resInfo)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write(resJSON)
		}

	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// startP2PServer start a server to received peer
func startP2PServer() error {
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/download", handleDownload)
	err := serv.ListenAndServe()
	return err
}

// closeP2PServer close the running P2PServer
func closeP2PServer() error {
	err := serv.Close()
	return err
}
