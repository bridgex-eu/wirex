package write

import (
	"net/http"

	"github.com/bridgex-eu/wirex"
)

type WriterWithHeaders struct {
	writer  wirex.Writer
	headers []wirex.HeaderWriter
}

func (m *WriterWithHeaders) WriteResponse(w http.ResponseWriter, r *http.Request) {
	for _, item := range m.headers {
		if err := item.WriteHeader(w, r); err != nil {
			err.WriteResponse(w, r)
			return
		}
	}

	m.writer.WriteResponse(w, r)
}

func WithHeader(writer wirex.Writer, headers ...wirex.HeaderWriter) *WriterWithHeaders {
	return &WriterWithHeaders{
		writer:  writer,
		headers: headers,
	}
}
