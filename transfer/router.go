package transfer

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	respCodeSucces      = iota + 200 // successful code
	respCodePathMissing              // miss path
	respCodeNotFound                 // path not found
)

// server response JSON schema
type responseJSON struct {
	Code         int    `json:"code"`                    // server action code
	Data         any    `json:"data,omitempty"`          // server returned data
	ErrorMessage string `json:"error_message,omitempty"` // server returned error when code is not respCodeSucces(200)
}

// send JSON string error
func sendError(errMsg string, errCode int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	resInfo := responseJSON{
		errCode,
		nil,
		errMsg,
	}
	resJSON, err := json.Marshal(resInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(resJSON)
}

// handleList handle /list route.
//
// it responses send file list JSON.
func handleList(w http.ResponseWriter, req *http.Request, files_list []TransferFile) {
	method := req.Method
	switch method {
	case http.MethodGet:
		resInfo := responseJSON{respCodeSucces, files_list, ""}
		resJSON, err := json.Marshal(resInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(resJSON)), 10))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(resJSON)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// handleDownload handle /download route.
//
// it responses byte stream according to route query path filed.
func handleDownload(w http.ResponseWriter, req *http.Request, files_list []TransferFile) {
	method := req.Method
	ranges, ok := req.Header["Range"]
	rangeH := make([]int64, 2)
	if ok {
		rangeH, _ = parseRangeHeader(ranges[0])
	}
	switch method {
	case http.MethodGet:
		path, ok := req.URL.Query()["path"]
		w.Header().Set("Access-Control-Allow-Headers", "*")
		// path filed lose
		if !ok || len(path) == 0 {
			sendError("query path filed is necessary.", respCodePathMissing, w)
			return
		}
		filePath := path[0]
		hasPathExisted := some(files_list, func(v TransferFile, all []TransferFile, i int) bool {
			return strings.Compare(v.Path, filePath) == 0
		})
		// path filed exist and path exist in file list
		if hasPathExisted {
			file, err := os.Open(filePath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			defer file.Close()
			fileInfo, _ := file.Stat()
			data, isChunked, err := readByte(file, rangeH[0])
			size := fileInfo.Size()
			if isChunked {
				w.Header().Set("Content-Range", formatRangeHeader(rangeH[1], rangeH[1]+MAX_BUFFER_BYTE, size))
			} else {
				w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
			}
			w.Write(data)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			// path not exist in file list
			sendError("query path has been deleted by send peer.", respCodeNotFound, w)
		}

	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}
