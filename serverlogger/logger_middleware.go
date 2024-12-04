package serverlogger

import (
	"datcha/servercommon"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func LoggerWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1<<16)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				slog.Error(fmt.Sprintf("Panic. Error= %v\n %s", err, buf))
				http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
			}
		}()
		writter := NotFoundResponseWriter{
			ResponseWriter: w,
		}
		sttime := time.Now()
		slog.Log(r.Context(), LevelRequest, "Request",
			slog.String("method", r.Method),
			slog.String("uri", r.RequestURI),
			slog.String("remote", r.RemoteAddr))
		h.ServeHTTP(&writter, r)
		slog.Log(r.Context(), LevelRequest, "Request done",
			slog.String("interval", strconv.FormatInt(time.Since(sttime).Milliseconds(), 10)))
	})
}
