package providers

import (
	"bytes"
	"encoding/xml"
)

func EscapeXML(d string) string {
	buf := &bytes.Buffer{}
	xml.Escape(buf, []byte(d))
	return buf.String()
}
