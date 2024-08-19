package http

import (
	"log/slog"
	"net/http"
)

func ServeHTTP(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request", "method", r.Method, "url", r.URL.Path, "remote", r.RemoteAddr, "user-agent", r.UserAgent(), "proto", r.Proto)

		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// return the elo.svg file if the request is for /karting
		if r.URL.Path == "/karting" {
			http.ServeFile(w, r, "elo.svg")
			return
		}

		w.WriteHeader(http.StatusNotFound)

		_, err := w.Write([]byte("Not found"))
		if err != nil {
			slog.Error("failed to write response", "error", err)
		}
	})

	slog.Info("Starting HTTP server", "port", port)
	http.ListenAndServe(":8080", nil)
}
