package write

import (
	"log/slog"
	"net/http"

	"github.com/bridgex-eu/wirex"
)

type BlobData struct {
	Status      int
	ContentType string
	Data        []byte
}

var _ wirex.Writer = &BlobData{}

func (b *BlobData) WriteResponse(w http.ResponseWriter, r *http.Request) {
	writeHeader := w.Header()
	if writeHeader.Get(wirex.HeaderContentType) == "" {
		writeHeader.Set(wirex.HeaderContentType, b.ContentType)
	}

	w.WriteHeader(b.Status)

	_, err := w.Write(b.Data)
	if err != nil {
		slog.Error("error occurred while trying to write to the response writer", "error", err)
	}
}

func Blob(status int, contentType string, data []byte, other ...wirex.HeaderWriter) wirex.Writer {
	return WithHeader(&BlobData{status, contentType, data}, other...)
}
