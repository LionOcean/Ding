package transfer

import (
	"net/http"
)

func startP2PServer() {
	http.HandleFunc("/list", func(w http.ResponseWriter, req *http.Request) {

	})

}
