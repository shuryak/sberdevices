package util

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/shuryak/sberhack/pkg/smarthome/endpoint"
)

func RunEndpoint(_ context.Context, client *http.Client, e *endpoint.Endpoint, dest ...interface{}) (int, error) {
	req, err := http.NewRequest(e.Method, e.PreparedURL(), bytes.NewReader(e.Body))
	if err != nil {
		return 0, err
	}

	if e.Headers != nil {
		req.Header = e.Headers
	}

	if len(e.Headers.Get("Accept")) == 0 {
		req.Header.Set("Accept", "application/json, text/plain, */*")
	}
	if len(e.Headers.Get("Content-Type")) == 0 {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	// Just ignore error
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	if dest != nil {
		var body []byte

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return resp.StatusCode, err
		}

		for i := 0; i < len(dest); i++ {
			err = json.Unmarshal(body, &dest[i])
		}
	}

	return resp.StatusCode, err
}
