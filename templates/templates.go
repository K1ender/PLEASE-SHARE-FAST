package templates

import (
	"embed"
	"fmt"
)

//go:embed files/*
var fs embed.FS

type template string

var (
	Upload template = "upload.html"
)

func GetTemplate(name template) ([]byte, error) {
	return fs.ReadFile(fmt.Sprintf("files/%s", name))
}
