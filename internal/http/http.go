package http

import (
	"fmt"
	"log/slog"
	"net/http"
)

func StartHttpServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Device successfully provisioned")
	})

	slog.Info("HTTP server on :80 /setup")

	if err := http.ListenAndServe(":80", mux); err != nil {
		slog.Error("http server failed", "error", err)
	}
}
