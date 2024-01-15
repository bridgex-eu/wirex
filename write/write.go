package write

import (
	"net/http"

	"github.com/bridgex-eu/wirex"
)

func setContentType(w http.ResponseWriter, value string) {
	writeHeader := w.Header()
	if writeHeader.Get(wirex.HeaderContentType) == "" {
		writeHeader.Set(wirex.HeaderContentType, value)
	}
}
