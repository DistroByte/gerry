package http

import (
	"log/slog"
	"net/http"
)

func ServeHTTP(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", "method", r.Method, "url", r.URL.Path, "ip", r.Header.Get("Cf-Connecting-Ip"), "user-agent", r.UserAgent())

		// ensure we only respond to GET requests
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// return a 200 OK for /health
		if r.URL.Path == "/health" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("OK"))
			if err != nil {
				slog.Error("failed to write response", "error", err)
			}
			return
		}

		if r.URL.Path == "/karting" {
			w.Header().Set("Content-Type", "image/svg+xml")
			http.ServeFile(w, r, "elo.svg")
			return
		}

		if r.URL.Path == "/karting.png" {
			w.Header().Set("Content-Type", "image/png")
			http.ServeFile(w, r, "elo.png")
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
