package http

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func initHTTPServer() {
	c := alice.New()

	c = c.Append(hlog.NewHandler(log.Logger.With().Str("component", "http").Logger()))

	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Str("Cf-Connecting-Ip", r.Header.Get("Cf-Connecting-Ip")).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user-agent"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))
	c = c.Append(hlog.RefererHandler("referer"))

	// create a new mux
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /karting", kartingHandler)
	mux.HandleFunc("GET /karting.png", kartingPNGHandler)
	mux.HandleFunc("GET /*", notFoundHandler)
	mux.HandleFunc("/*", methodNotAllowedHandler)

	http.Handle("/", c.Then(mux))
}

func ServeHTTP() {
	initHTTPServer()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Msg("")
	}
}

func kartingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	http.ServeFile(w, r, "elo.svg")
}

func kartingPNGHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, "elo.png")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("Not found"))
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Msg("")
	}
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
