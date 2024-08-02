package debugsrv

import (
	"context"
	"net/http"
	"os"
	"route256/common/pkg/env"
)

const (
	EnvAppDebugSrv         = "APP_DEBUG_SRV"
	EnvAppDebugSrvHTTPAddr = "APP_DEBUG_SRV_HTTP_ADDR"
)

func Run(ctx context.Context, opts ...func(mux *http.ServeMux)) {
	if !isEnabled() {
		return
	}

	conf := getConfig(opts...)
	server := NewServer(conf)
	server.Run(ctx)
}

func getConfig(opts ...func(mux *http.ServeMux)) Config {
	return Config{
		HTTPAddr: env.GetEnvVar(EnvAppDebugSrvHTTPAddr),
		options:  opts,
	}
}

func isEnabled() bool {
	v, ok := os.LookupEnv(EnvAppDebugSrv)
	if !ok {
		return true
	}

	return v != "disabled"
}
