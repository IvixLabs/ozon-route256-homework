package pprof

import (
	"net/http"
	"net/http/pprof"
)

func WithPprof() func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {

		prefix := "/debug/pprof"
		mux.HandleFunc(prefix+"/*", pprof.Index)
		mux.HandleFunc(prefix+"/cmdline", pprof.Cmdline)
		mux.HandleFunc(prefix+"/profile", pprof.Profile)
		mux.HandleFunc(prefix+"/symbol", pprof.Symbol)
		mux.HandleFunc(prefix+"/trace", pprof.Trace)
	}
}
