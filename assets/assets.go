package assets

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
)

//go:generate go-bindata -ignore "\\.go" -pkg assets -prefix dist -o bindata.go ./dist/...
//go:generate go fmt bindata.go
//go:generate sed -i.bak "s/Css/CSS/" bindata.go
//go:generate rm bindata.go.bak

// Load initializes the static files.
func Load() http.FileSystem {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "",
	}
}
