package httphelpers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/chromedp/cdproto/har"
)

type harLoggerStorageCtxKeyType string

var HarLoggerStorageCtxKey harLoggerStorageCtxKeyType = "harLogger"

type HarLoggerStorage struct {
	HarDataLock sync.Mutex
	HarData     *har.HAR
}

func (s *HarLoggerStorage) GetJSONData() ([]byte, error) {
	s.HarDataLock.Lock()
	defer s.HarDataLock.Unlock()
	return json.Marshal(s.HarData)
}

func WithHarLoggerStorage(ctx context.Context) context.Context {
	return context.WithValue(ctx, HarLoggerStorageCtxKey, &HarLoggerStorage{
		HarData: &har.HAR{
			Log: &har.Log{
				Version: "1.2",
				Creator: &har.Creator{
					Name:    "paczkobot",
					Version: "unk",
					Comment: "Generated with HarLoggerStorage",
				},
				Entries: []*har.Entry{},
			},
		},
	})
}

func GetHarLoggerStorage(ctx context.Context) *HarLoggerStorage {
	if ctx.Value(HarLoggerStorageCtxKey) == nil {
		return nil
	}
	return ctx.Value(HarLoggerStorageCtxKey).(*HarLoggerStorage)
}

type HarLoggerTransport struct {
	Transport http.RoundTripper
}

func NewHarLoggerTransport(rt http.RoundTripper) *HarLoggerTransport {
	return &HarLoggerTransport{rt}
}

func (t *HarLoggerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	storage := GetHarLoggerStorage(req.Context())
	if storage == nil {
		return t.Transport.RoundTrip(req)
	}
	var requestBody bytes.Buffer
	if req.Body != nil {
		requestBody = bytes.Buffer{}
		// read whole body
		_, err := requestBody.ReadFrom(req.Body)
		if err != nil {
			return nil, err
		}
		// restore body
		req.Body = io.NopCloser(bytes.NewReader(requestBody.Bytes()))
	}
	entry := &har.Entry{
		StartedDateTime: time.Now().Format(time.RFC3339),
		Request: &har.Request{
			Method:      req.Method,
			URL:         req.URL.String(),
			HTTPVersion: req.Proto,
			Headers:     mapHeadersToHar(req.Header),
			Cookies:     mapCookiesToHar(req.Cookies()),
			QueryString: []*har.NameValuePair{},
			BodySize:    -1,
			HeadersSize: -1,
		},
	}

	if entry.Request.HTTPVersion == "" {
		entry.Request.HTTPVersion = "HTTP/1.1"
	}

	if req.Body != nil {
		entry.Request.BodySize = int64(requestBody.Len())
		entry.Request.PostData = &har.PostData{
			MimeType: req.Header.Get("Content-Type"),
			Text:     requestBody.String(),
		}
	}

	resp, err := t.Transport.RoundTrip(req)

	// defer is used to make sure that we save the request even if there is an error
	defer func() {
		storage.HarDataLock.Lock()
		defer storage.HarDataLock.Unlock()
		storage.HarData.Log.Entries = append(storage.HarData.Log.Entries, entry)
	}()

	if err != nil {
		return nil, err
	}

	entry.Response = &har.Response{
		Status:      int64(resp.StatusCode),
		StatusText:  resp.Status,
		HTTPVersion: resp.Proto,
		Cookies:     mapCookiesToHar(resp.Cookies()),
		Headers:     mapHeadersToHar(resp.Header),
		Content:     &har.Content{},
		BodySize:    -1,
		HeadersSize: -1,
		RedirectURL: resp.Header.Get("Location"),
	}

	if entry.Response.HTTPVersion == "" {
		entry.Response.HTTPVersion = "HTTP/1.1"
	}

	respBody := bytes.Buffer{}
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(respBody.Bytes()))

	encodedBody := base64.StdEncoding.EncodeToString(respBody.Bytes())

	entry.Response.BodySize = int64(respBody.Len())
	entry.Response.Content = &har.Content{
		Size:     int64(respBody.Len()),
		MimeType: resp.Header.Get("Content-Type"),
		Text:     encodedBody,
		Encoding: "base64",
	}

	return resp, nil
}

func NewClientWithLogger() *http.Client {
	return &http.Client{
		Transport: NewHarLoggerTransport(http.DefaultTransport),
	}
}

func mapCookiesToHar(inp []*http.Cookie) []*har.Cookie {
	cookies := []*har.Cookie{}
	for _, c := range inp {
		cookies = append(cookies, &har.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Path:     c.Path,
			Domain:   c.Domain,
			Expires:  c.Expires.Format(time.RFC3339),
			HTTPOnly: c.HttpOnly,
			Secure:   c.Secure,
		})
	}
	return cookies
}

func mapHeadersToHar(inp http.Header) []*har.NameValuePair {
	headers := []*har.NameValuePair{}
	for k, v := range inp {
		for _, vv := range v {
			headers = append(headers, &har.NameValuePair{
				Name:  k,
				Value: vv,
			})
		}
	}
	return headers
}
