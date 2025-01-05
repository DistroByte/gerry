package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/distrobyte/gerry/internal/config"
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
	mux.HandleFunc("GET /karting", kartingHTMLHandler)
	mux.HandleFunc("GET /karting.svg", kartingSVGHandler)
	mux.HandleFunc("GET /karting.png", kartingPNGHandler)
	mux.HandleFunc("GET /karting.json", kartingJSONHandler)
	mux.HandleFunc("GET /*", notFoundHandler)
	mux.HandleFunc("/*", methodNotAllowedHandler)

	http.Handle("/", c.Then(mux))

	log.Info().Msg("http server initialized. listening on port " + fmt.Sprintf("%d", config.GetHTTPPort()))
}

func ServeHTTP() {
	initHTTPServer()

	err := http.ListenAndServe(fmt.Sprintf(":%d", config.GetHTTPPort()), nil)
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

func kartingHTMLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	htmlContent := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Karting SVG</title>
        <script defer data-domain="gerry.dbyte.xyz" src="https://plausible.dbyte.xyz/js/script.js"></script>
    </head>
    <body>
        <object type="image/svg+xml" data="karting.svg"></object>
    </body>
    </html>`

	_, err := w.Write([]byte(htmlContent))
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Msg("")
		return
	}
}

func kartingSVGHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	http.ServeFile(w, r, "elo.svg")
}

func kartingPNGHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, "elo.png")
}

func kartingJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "karting.json")
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
