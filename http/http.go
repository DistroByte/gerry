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

	r.Get("/health", healthHandler)

	// Serve assets directory from root
	// This allows direct access to files like /elo.html, /elo.png, /karting.json, etc.
	r.Handle("/*", http.FileServer(http.Dir("assets")))

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
