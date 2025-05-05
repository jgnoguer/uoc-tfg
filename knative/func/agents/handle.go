package function

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
)

// Handle an HTTP Request.
func Handle(w http.ResponseWriter, r *http.Request) {
	/*
	 * YOUR CODE HERE
	 *
	 * Try running `go test`.  Add more test as you code in `handle_test.go`.
	 */

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Received request")
	fmt.Printf("%q\n", dump)
	fmt.Fprintf(w, "%q", dump)

	dirError := os.Mkdir("/data/hola/", os.ModeDir)
	if dirError != nil {
		slog.Error("Fail to create directory" + dirError.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
