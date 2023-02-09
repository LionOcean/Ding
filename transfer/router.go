package transfer

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	responseCodeSuccess = iota + 200 // successful code
	responseCodeFail                 // faliure code
)

// TransferFile is the file that send to received peer
type TransferFile struct {
	Path string `json:"path"` // file absolute path
	Size int    `json:"size"` // file byte length
}

// server response JSON schema
type responseJSON struct {
	Code         int    `json:"code"`          // server action code
	Data         any    `json:"data"`          // server returned data
	ErrorMessage string `json:"error_message"` // server returned error when code is not responseCodeSuccess(200)
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
		hasPathExisted := some(files_list, func(v TransferFile, all []TransferFile, i int) bool {
			return strings.Compare(v.Path, filePath) == 0
		})
		// path filed exist and path exist in file list
		if ok && hasPathExisted {
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
			if !hasPathExisted {
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
