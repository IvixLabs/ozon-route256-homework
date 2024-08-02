package swagger

import (
	"embed"
	"io"
	"io/fs"
	"log"
)

//go:embed assets
var assets embed.FS

func GetFileData() []byte {

	resourceDir, err := fs.Sub(assets, "assets")
	file, err := resourceDir.Open("loms.swagger.json")
	if err != nil {
		log.Panicln(err)
	}

	rawFile, err := io.ReadAll(file)
	if err != nil {
		log.Panicln(err)
	}

	return rawFile
}
