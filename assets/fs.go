package assets

import (
	"embed"
)

//go:embed certs/*
//go:embed static/*
//go:embed views/*
var FS embed.FS
