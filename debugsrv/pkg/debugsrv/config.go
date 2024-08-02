package debugsrv

import (
	"net/http"
)

type Config struct {
	HTTPAddr string
	options  []func(*http.ServeMux)
}
