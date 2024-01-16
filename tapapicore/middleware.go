package tapapicore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"crud/log"
	"crud/tapcontext"
	"github.com/felixge/httpsnoop"
	"github.com/google/uuid"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
)

// TraceID is the header that is passed in outgoing response
const TraceID = "traceid"

// SetTraceID is a http middleware that checks for "Traceparent" header in the
// incoming requests. If present, parse and extract the traceid. Else create a custom traceid
// and sets it as new "Traceparent" header in incoming request.
// Also sets that traceid in outgoing response header.
func SetTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var traceID apm.TraceID
		if values := r.Header[apmhttp.W3CTraceparentHeader]; len(values) == 1 && values[0] != "" {
			if c, err := apmhttp.ParseTraceparentHeader(values[0]); err == nil {
				traceID = c.Trace
			}
		}
		if err := traceID.Validate(); err != nil {
			uuidId := uuid.New()
			var spanID apm.SpanID
			var traceOptions apm.TraceOptions
			copy(traceID[:], uuidId[:])
			copy(spanID[:], traceID[8:])
			traceContext := apm.TraceContext{
				Trace:   traceID,
				Span:    spanID,
				Options: traceOptions.WithRecorded(true),
			}
			r.Header.Set(apmhttp.W3CTraceparentHeader, apmhttp.FormatTraceparentHeader(traceContext))
		}

		w.Header().Set(TraceID, traceID.String())
		r.Header.Set(requestID, traceID.String())
		next.ServeHTTP(w, r)
	})
}

// middleware to handle panics
func recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			ctx := tapcontext.UpgradeCtx(r.Context())
			rec := recover()
			if rec != nil {
				span, _ := apm.StartSpan(ctx.Context, "recovery", "custom")
				defer span.End()
				trace := string(debug.Stack())
				trace = strings.Replace(trace, "\n", "    ", -1)
				trace = strings.Replace(trace, "\t", "    ", -1)
				log.GenericError(ctx, fmt.Errorf("%v", rec),
					log.FieldsMap{
						"msg":        "recovering from panic",
						"stackTrace": trace,
					})
				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error",
				})
				e := apm.DefaultTracer.Recovered(rec)
				e.SetSpan(span)
				e.Send()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

// enableCorsMiddleware : sets the CORS for this service
func enableCorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}

}

// enableLogging : sets the logging for this service
func enableLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := tapcontext.UpgradeCtx(r.Context())
		// Avoid logging ping API
		rawBody, _ := ioutil.ReadAll(r.Body)
		if len(rawBody) > 0 {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
		}

		if strings.Contains(r.RequestURI, "ping") {
			return
		}
		m := log.FieldsMap{
			"method":  r.Method,
			"url":     r.RequestURI,
			"reqBody": string(rawBody),
		}
		// Log only the request type and URI
		log.GenericInfo(ctx, "trace", m)
		next.ServeHTTP(w, r)
	}
}

// createContext : sets a custom context for this service
func createContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header
		ctx := r.Context()
		reqID := header.Get(requestID)
		if reqID == "" {
			reqID = strings.ReplaceAll(uuid.NewString(), "-", "")
		}
		email, token, app :=
			header.Get(userEmail), header.Get(tapApiToken), header.Get(application)
		locale := header.Get(locale)

		tapCtx := tapcontext.TapContext{
			RequestID:   reqID,
			UserEmail:   email,
			TapApiToken: token,
			Application: app,
			Locale:      locale,
		}

		ctx = tapcontext.WithTapCtx(ctx, tapCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// logRequest logs each HTTP incoming Requests
func logRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m := httpsnoop.CaptureMetrics(handler, w, r)
		log.HTTPLog(constructHTTPLog(r, m, time.Since(start)))
	}
}

// constructHTTPLogMessage
func constructHTTPLog(r *http.Request, m httpsnoop.Metrics, duration time.Duration) string {
	ctx := r.Context().Value(tapcontext.TAPCtx)

	rawBody, _ := ioutil.ReadAll(r.Body)
	if len(rawBody) > 0 {
		r.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
	}
	if ctx != nil {
		tCtx := ctx.(tapcontext.TapContext)
		return fmt.Sprintf("|%s|%s|%s|%s|%s|%d|%d|%s|%s|%s|",
			// Cannot modify original request/obtain apiContext through gorilla context, hence we won't get the apiContext data from the request object.
			tCtx.UserEmail+":",
			"requestId="+tCtx.RequestID,
			r.RemoteAddr,
			r.Method,
			r.URL,
			m.Code,
			m.Written,
			r.UserAgent(),
			duration,
			"Body:"+string(rawBody),
		)
	}
	return fmt.Sprintf("|%s|%s|%s|%d|%d|%s|%s|%s|",
		// Cannot modify original request/obtain apiContext through gorilla context, hence we won't get the apiContext data from the request object.
		r.RemoteAddr,
		r.Method,
		r.URL,
		m.Code,
		m.Written,
		r.UserAgent(),
		duration,
		"Body:"+string(rawBody),
	)

}

// enableGzip : zips the response if Accept-Encoding is gzip
func enableCompression(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "br") && !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next(w, r)
			return
		} else if !strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			gzr := pool.Get().(*gzipResponseWriter)
			gzr.statusCode = 0
			gzr.headerWritten = false
			gzr.ResponseWriter = w
			gzr.w.Reset(w)
			defer func() {
				// gzr.w.Close will write a footer even if no data has been written.
				// StatusNotModified and StatusNoContent expect an empty body so don't close it.
				if gzr.statusCode != http.StatusNotModified && gzr.statusCode != http.StatusNoContent {
					if err := gzr.w.Close(); err != nil {
						ctx := tapcontext.UpgradeCtx(r.Context())
						log.GenericError(ctx, err, nil)
					}
				}
				pool.Put(gzr)
			}()
			next(gzr, r)
			return
		}
		br := poolbr.Get().(*brotliResponseWriter)
		br.statusCode = 0
		br.headerWritten = false
		br.ResponseWriter = w
		br.w.Reset(w)

		defer func() {
			// brotli.w.Close will write a footer even if no data has been written.x
			// StatusNotModified and StatusNoContent expect an empty body so don't close it.
			if br.statusCode != http.StatusNotModified && br.statusCode != http.StatusNoContent {
				if err := br.w.Close(); err != nil {
					ctx := tapcontext.UpgradeCtx(r.Context())
					log.GenericError(ctx, err, nil)
				}
			}
			poolbr.Put(br)
		}()
		next(br, r)
	}
}

func enableUserValidationForApplication(inner http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span, _ := apm.StartSpan(ctx, "UserValidationForApplication", "authorisation")
		tCtx := new(tapcontext.TContext)
		tCtx.TapContext = r.Context().Value(tapcontext.TAPCtx).(tapcontext.TapContext)
		if tCtx.UserEmail == "" {
			w.Write([]byte("Basic values should be present in header"))
			return
		}
		ctx = tapcontext.WithTapCtx(ctx, tCtx.TapContext)
		span.End()
		inner.ServeHTTP(w, r.WithContext(ctx))
	}
}
