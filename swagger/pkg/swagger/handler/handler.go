package handler

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
)

//go:embed assets
var assets embed.FS

func NewFileHandler() http.Handler {
	resourceDir, err := fs.Sub(assets, "assets")
	if err != nil {
		panic(err)
	}

	resourcesFS := http.FS(resourceDir)

	return http.FileServer(resourcesFS)
}

func NewApiDocsHandler(restAddr string, swaggerJson []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var jsonObj map[string]interface{}

		err := json.Unmarshal(swaggerJson, &jsonObj)
		if err != nil {
			log.Fatalln(err)
		}

		jsonObj["host"] = restAddr

		newRaw, err := json.Marshal(jsonObj)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = w.Write(newRaw)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
