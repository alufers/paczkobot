package commondata

import "net/http"

func SetCommonHTTPHeaders(header *http.Header) {

	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")
	header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	header.Set("Accept-Language", "en-US,en;q=0.5")
	header.Set("Upgrade-Insecure-Requests", "1")
	header.Set("Pragma", "no-cache")
	header.Set("Cache-Control", "no-cache")
}
