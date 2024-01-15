package from

import (
	"net/http"

	"github.com/bridgex-eu/wirex"
)

func Bind(r *http.Request, items ...wirex.FromRequest) wirex.HTTPError {
	for _, item := range items {
		if err := item.FromRequest(r); err != nil {
			return err
		}
	}

	return nil
}
