package doc

import (
	"embed"
	"net/http"
	"strings"
)

//go:embed api/*
//go:embed events/*
var fs embed.FS

type prefixTrimmingFS struct {
	http.FileSystem

	prefix string
}

func (p *prefixTrimmingFS) Open(name string) (http.File, error) {
	f, err := p.FileSystem.Open(strings.TrimPrefix(name, p.prefix))
	return f, err
}

var Handler = http.FileServer(&prefixTrimmingFS{
	FileSystem: http.FS(fs),
	prefix:     "/doc",
})
