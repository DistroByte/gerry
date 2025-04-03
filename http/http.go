package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/distrobyte/gerry/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func initHTTPServer() {

	r := chi.NewRouter()
	r.Use(requestIDMiddleware)
	r.Use(zerologMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// use the kartingHTML handler to serve the main page
		kartingHTMLHandler(w, r)
	})

	r.Get("/health", healthHandler)
	r.Get("/karting", kartingHTMLHandler)
	r.Get("/karting.svg", kartingSVGHandler)
	r.Get("/karting.png", kartingPNGHandler)
	r.Get("/karting.json", kartingJSONHandler)
	r.NotFound(notFoundHandler)
	r.MethodNotAllowed(methodNotAllowedHandler)

	log.Info().Msgf("Starting server on port %d", config.GetHTTPPort())
	log.Fatal().Err(http.ListenAndServe(fmt.Sprintf(":%d", config.GetHTTPPort()), r)).Msg("")
	log.Info().Msgf("Server started on port %d", config.GetHTTPPort())
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		logger := log.With().Str("request_id", requestID).Logger()
		r = r.WithContext(logger.WithContext(r.Context()))
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}

func zerologMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		logger := zerolog.Ctx(r.Context())
		defer func() {
			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", ww.Status()).
				Int("bytes", ww.BytesWritten()).
				Str("remote", r.RemoteAddr).
				Dur("duration", time.Since(start)).
				Msg("handled request")
		}()

		next.ServeHTTP(ww, r)
	})
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
	if r.Method != http.MethodGet {
		methodNotAllowedHandler(w, r)
		return
	}

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
