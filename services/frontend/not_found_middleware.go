package frontend

import (
	"datcha/serverlogger"
	"fmt"
	"log/slog"
	"net/http"
	"path"
)

func (server FrontendService) fsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respWriter := &serverlogger.NotFoundResponseWriter{ResponseWriter: w}
		h.ServeHTTP(respWriter, r)
		if respWriter.Status == http.StatusNotFound {
			if server.NotFoundRedirectPage != "" {
				slog.Info(fmt.Sprintf("Page not found. Redirecting %s to %s", r.RequestURI, server.NotFoundRedirectPage))
				http.Redirect(w, r, server.NotFoundRedirectPage, http.StatusFound)
			} else if server.NotFoundFile != "" {
				fileName := path.Join(server.FrontendFolder, server.NotFoundFile)
				slog.Info(fmt.Sprintf("Page not found. Send file %s", fileName))
				// Previously Content-Type header may be set by FileServer.
				// It will not updated by ServeFile if it is already set. So remove it
				w.Header().Del("Content-Type")
				http.ServeFile(w, r, fileName)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	})
}
