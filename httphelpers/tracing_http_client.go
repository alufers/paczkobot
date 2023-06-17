package httphelpers

import (
	"log"
	"net/http"
	"net/http/httptrace"
)

type TracingHttpClient struct {
	HttpClient Client
}

func NewTracingHttpClient(httpClient Client) *TracingHttpClient {
	return &TracingHttpClient{
		HttpClient: httpClient,
	}
}

func (c *TracingHttpClient) Do(req *http.Request) (*http.Response, error) {
	log.Printf("Sending %s request to %s", req.Method, req.URL)
	trace := &httptrace.ClientTrace{
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			log.Printf("DNS Info: %+v\n", dnsInfo)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			log.Printf("Got Conn: %+v\n", connInfo)
		},
		ConnectStart: func(network, addr string) {
			log.Printf("Connect Start: %s %s\n", network, addr)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	return http.DefaultTransport.RoundTrip(req)
}
