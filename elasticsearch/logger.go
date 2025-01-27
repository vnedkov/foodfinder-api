package elasticsearch

import (
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// https://github.com/elastic/go-elasticsearch/blob/main/_examples/logging/custom.go

// CustomLogger implements the elastictransport.Logger interface.
type CustomLogger struct {
	zerolog.Logger
}

// LogRoundTrip prints the information about request and response.
func (l *CustomLogger) LogRoundTrip(
	req *http.Request,
	res *http.Response,
	err error,
	start time.Time,
	dur time.Duration,
) error {
	var (
		e    *zerolog.Event
		nReq int64
		nRes int64
	)

	// Set error level.
	switch {
	case err != nil:
		e = l.Error()
	case res != nil && res.StatusCode > 0 && res.StatusCode < 300:
		e = l.Info()
	case res != nil && res.StatusCode > 299 && res.StatusCode < 500:
		e = l.Warn()
	case res != nil && res.StatusCode > 499:
		e = l.Error()
	default:
		e = l.Error()
	}

	// Count number of bytes in request and response.
	var reqBody []byte
	var resBody []byte
	if req != nil && req.Body != nil && req.Body != http.NoBody {
		reqBody, _ = io.ReadAll(req.Body)
		nReq = int64(len(reqBody))
	}
	if res != nil && res.Body != nil && res.Body != http.NoBody {
		resBody, _ = io.ReadAll(res.Body)
		nRes = int64(len(resBody))
	}

	// Log event.
	e.Str("method", req.Method).
		Int("status_code", res.StatusCode).
		Dur("duration", dur).
		Int64("req_bytes", nReq).
		Int64("res_bytes", nRes).
		Str("request_body", string(reqBody)).
		Str("response_body", string(resBody)).
		Msg(req.URL.String())

	return nil
}

// RequestBodyEnabled makes the client pass request body to logger
func (l *CustomLogger) RequestBodyEnabled() bool { return true }

// RequestBodyEnabled makes the client pass response body to logger
func (l *CustomLogger) ResponseBodyEnabled() bool { return true }
