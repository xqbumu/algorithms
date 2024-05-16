package main

import (
	"algorithms/examples/katana/pkg/parser"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
)

func main() {
	// Headless runs the browser on foreground, you can also use flag "-rod=show"
	// Devtools opens the tab in each new tab opened automatically
	l := launcher.New().
		Set("disable-blink-features", "AutomationControlled").
		Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36").
		UserDataDir("userdata").
		Headless(false).
		Devtools(true)

	defer l.Cleanup() // remove launcher.FlagUserDataDir

	url := l.MustLaunch()

	// Trace shows verbose debug information for each action executed
	// Slowmotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	browser := rod.New().
		ControlURL(url).
		Trace(true).
		SlowMotion(2 * time.Second).
		MustConnect()

	// ServeMonitor plays screenshots of each tab. This feature is extremely
	// useful when debugging with headless mode.
	// You can also enable it with flag "-rod=monitor"
	launcher.Open(browser.ServeMonitor(""))

	defer browser.MustClose()

	url = "https://www.baidu.com"
	// url = "https://cdek-express.com/?tracking=1535486778"
	// url = "https://www.sibtrans.ru/status/?n=57iTUiG7X2"
	page := browser.NoDefaultDevice().MustPage(url)

	result, err := parser.CdekExpress(page)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	utils.Pause() // pause goroutine
}
