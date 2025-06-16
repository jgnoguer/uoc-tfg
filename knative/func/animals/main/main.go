package main

import (
	"errors"
	"function"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func KnativeHandler(w http.ResponseWriter, r *http.Request) {
	function.Handle(w, r)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	server := http.NewServeMux()

	server.HandleFunc("/", KnativeHandler)

	slog.Info("starting server")

	err := http.ListenAndServe(":8080", server)

	if errors.Is(err, http.ErrServerClosed) {
		slog.Info("server closed")
	} else if err != nil {
		slog.Error("error starting server: %s\n", err.Error())
		os.Exit(1)
	}
}
