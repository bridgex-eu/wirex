package write

import "github.com/bridgex-eu/wirex"

func String(status int, value string, other ...wirex.HeaderWriter) wirex.Writer {
	return Blob(status, wirex.MIMETextPlainCharsetUTF8, []byte(value))
}
