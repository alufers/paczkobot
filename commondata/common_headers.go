package commondata

import "net/http"

func SetCommonHTTPHeaders(header *http.Header) {

	header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0")
	header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	header.Set("Accept-Language", "en-US,en;q=0.5")
	header.Set("Upgrade-Insecure-Requests", "1")
	header.Set("Pragma", "no-cache")
	header.Set("Cache-Control", "no-cache")
}
