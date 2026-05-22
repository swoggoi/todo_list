package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/swoggoi/todo_list/internal/core/logger"
	core_http_responce "github.com/swoggoi/todo_list/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func RequestID() Middleware {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_responce.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incomind HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)
			next.ServeHTTP(rw, r)
			log.Debug(
				"<<< done HTTP request",
				zap.Int("status code", rw.GetStatusCode()),
				zap.Duration("latency", time.Now().Sub(before)),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				ctx := r.Context()
				log := core_logger.FromContext(ctx)
				responseHandler := core_http_responce.NewHTTPResponseHandler(log, w)
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"during handle http request got unexpected panic",
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
