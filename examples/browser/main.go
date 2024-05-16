package main

import (
	"log"
	"math/rand"

	"github.com/go-rod/rod/lib/proto"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/zlsgo/browser"
)

func main() {
	b, err := browser.New(func(o *browser.Options) {
		o.AutoKill = false
		o.Debug = true
		o.Devtools = true
		o.Incognito = true
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		b.Close()
	}()

	err = b.Open("https://github.com/arkjit", func(p *browser.Page) error {
		zlog.Info(p.MustInfo().Title)
		el := p.MustElementX("/html/body/div[1]/div[1]/header/div/div[2]/div/nav/ul/li[4]/a")
		box := el.MustShape().Box()
		p.Mouse.MoveLinear(proto.Point{
			X: box.X + rand.Float64()*box.Height,
			Y: box.Y + rand.Float64()*box.Width,
		}, 10)
		log.Println(box)
		return nil
	}, func(o *browser.PageOptions) {
		// 劫持请求
		o.Hijack = map[string]browser.HijackProcess{
			"*": func(b *browser.Hijack) bool {
				// 屏蔽图片请求
				return !b.BlockImage()
			},
		}
	})
	if err != nil {
		panic(err)
	}

	wait := b.Browser.WaitEvent(proto.NetworkWebSocketClosed{})
	wait()
}
