package swagger

import (
	"net/http"
	"route256/common/pkg/env"
	"route256/swagger/pkg/swagger/handler"
)

const (
	EnvAppSwaggerRestAddr = "APP_SWAGGER_REST_ADDR"
)

func WithSwagger(swaggerJson []byte) func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {

		prefix := "/swagger"

		mux.Handle(prefix+"/*", http.StripPrefix("/swagger", handler.NewFileHandler()))

		restAddr := env.GetEnvVar(EnvAppSwaggerRestAddr)
		mux.HandleFunc(prefix+"/swagger.json", handler.NewApiDocsHandler(restAddr, swaggerJson))
	}
}
