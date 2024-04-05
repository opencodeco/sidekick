package proxy

import (
	"fmt"
	"io"
	"net/http"
)

func Proxy(appPort string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := http.Get("http://localhost:" + appPort + r.URL.Path)
		if err != nil {
			fmt.Println("error fetching", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("error reading body", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(res.StatusCode)
		w.Write(body)
	}
}
