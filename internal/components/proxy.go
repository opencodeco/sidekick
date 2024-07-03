package components

import (
	"io"
	"log/slog"
	"net/http"
)

func Proxy(appPort string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("proxying request", "path", r.URL.Path)

		res, err := http.Get("http://localhost:" + appPort + r.URL.Path)
		if err != nil {
			slog.Error("error fetching", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			slog.Error("error reading body", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(res.StatusCode)
		w.Write(body)
	}
}
