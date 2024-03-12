package routes

import (
	"net/http"

	"github.com/travisavey/baseline/app/logging"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.AccessLog.Info().Str("requst_url", r.RequestURI).Str("client", r.RemoteAddr).Str("user_agent", r.UserAgent()).Str("method", r.Method).Msg("server request")
		next.ServeHTTP(w, r)
	})
}
