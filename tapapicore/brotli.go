package tapapicore

import (
	"github.com/andybalholm/brotli"
	"net/http"
	"sync"
)

type brotliResponseWriter struct {
	http.ResponseWriter

	w             *brotli.Writer
	statusCode    int
	headerWritten bool
}

var (
	poolbr = sync.Pool{
		New: func() interface{} {
			w := brotli.NewWriterLevel(nil, brotli.BestSpeed)
			return &brotliResponseWriter{
				w: w,
			}
		},
	}
)

func (br *brotliResponseWriter) WriteHeader(statusCode int) {
	br.statusCode = statusCode
	br.headerWritten = true

	if br.statusCode != http.StatusNotModified && br.statusCode != http.StatusNoContent {
		br.ResponseWriter.Header().Del("Content-Length")
		br.ResponseWriter.Header().Set("Content-Encoding", "br")
	}

	br.ResponseWriter.WriteHeader(statusCode)
}
func (br *brotliResponseWriter) Write(b []byte) (int, error) {
	if _, ok := br.Header()["Content-Type"]; !ok {
		// If no content type, apply sniffing algorithm to un-gzipped body.
		br.ResponseWriter.Header().Set("Content-Type", http.DetectContentType(b))
	}

	if !br.headerWritten {
		// This is exactly what Go would also do if it hasn't been written yet.
		br.WriteHeader(http.StatusOK)
	}

	return br.w.Write(b)
}

func (br *brotliResponseWriter) Flush() {
	if br.w != nil {
		br.w.Flush()
	}

	if fw, ok := br.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}
