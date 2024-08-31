package main

import (
	"algorithms/examples/tadpole/logic"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/serialize/json"
	"github.com/urfave/cli"
)

//go:embed static/*
var staticfs embed.FS

func main() {
	app := cli.NewApp()

	app.Name = "tadpole"
	app.Author = "nano authors"
	app.Version = "0.0.1"
	app.Copyright = "nano authors reserved"
	app.Usage = "tadpole"

	// flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Value: ":8081",
			Usage: "game server address",
		},
	}

	app.Action = serve

	app.Run(os.Args)
}

func serve(cliCtx *cli.Context) error {
	components := &component.Components{}
	components.Register(logic.NewManager())
	components.Register(logic.NewWorld())

	// register all service
	options := []nano.Option{
		nano.WithIsWebsocket(true),
		nano.WithComponents(components),
		nano.WithSerializer(json.NewSerializer()),
		nano.WithCheckOriginFunc(func(_ *http.Request) bool { return true }),
		nano.WithWSPath("/ws"),
	}

	//nano.EnableDebug()
	log.SetFlags(log.LstdFlags | log.Llongfile)

	fSys, err := fs.Sub(staticfs, "static")
	if err != nil {
		panic(err)
	}

	http.Handle("/", http.FileServerFS(fSys))

	addr := cliCtx.String("addr")
	nano.Listen(addr, options...)
	return nil
}
