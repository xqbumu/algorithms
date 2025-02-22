package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
)

func main() {
	browser, clean := newBrowser()
	defer clean()

	RunSgpbusiness(browser)
}

func RunSgpbusiness(browser *rod.Browser) {
	for _, page := range browser.MustPages() {
		page.MustClose()
	}

	// page := stealth.MustPage(browser)

	// // new page clear cookie，inject stealth.JS
	// go browser.EachEvent(func(e *proto.TargetTargetCreated) {
	// 	if e.TargetInfo.Type != proto.TargetTargetInfoTypePage {
	// 		return
	// 	}
	// 	browser.MustPageFromTargetID(e.TargetInfo.TargetID).MustEvalOnNewDocument(stealth.JS)
	// })()

	// Create a new page
	// page := browser.MustPage("https://www.sgpbusiness.com/").MustWaitLoad()
	page := browser.MustPage("https://www.baidu.com/").MustWaitLoad()
	page.WaitStable(time.Millisecond * 800)

	search := page.MustElement("input[type=checkbox]")
	pt, err := search.WaitInteractable()
	if err != nil {
		panic(err)
	}
	page.Mouse.MoveTo(*pt)

	search.MustParent().Click(proto.InputMouseButtonLeft, 1)
	page.WaitStable(time.Millisecond * 800)

	search.MustInput("golang").MustType(input.Enter)
	page.WaitStable(time.Millisecond * 1200)

	slog.Info("Pause")
	utils.Pause()
}

func RunBaidu(browser *rod.Browser) {
	// Create a new page
	page := browser.MustPage("https://baidu.com").MustWaitLoad()
	page.WaitStable(time.Millisecond * 800)

	search := page.MustElement("#kw")
	pt, err := search.WaitInteractable()
	if err != nil {
		panic(err)
	}
	page.Mouse.MoveTo(*pt)

	search.MustParent().Click(proto.InputMouseButtonLeft, 1)
	page.WaitStable(time.Millisecond * 800)

	search.MustInput("golang").MustType(input.Enter)
	page.WaitStable(time.Millisecond * 1200)

	slog.Info("Pause")
	utils.Pause()
}

func RunGoogle(browser *rod.Browser) {
	// Create a new page
	page := browser.MustPage("https://google.com").MustWaitLoad()
	page.WaitStable(time.Millisecond * 800)

	search := page.MustElement("form textarea")
	pt, err := search.WaitInteractable()
	if err != nil {
		panic(err)
	}
	page.Mouse.MoveTo(*pt)

	search.MustParent().Click(proto.InputMouseButtonLeft, 1)
	page.WaitStable(time.Millisecond * 800)

	search.MustInput("golang").MustType(input.Enter)
	page.WaitStable(time.Millisecond * 1200)

	slog.Info("Pause")
	utils.Pause()
}

func RunDiscuz(browser *rod.Browser) {
	router := browser.HijackRequests()
	defer router.MustStop()

	seccodes := make(map[string]string, 0)

	router.MustAdd("*/misc.php?mod=seccode&update=*", func(ctx *rod.Hijack) {
		slog.Info("hijack", "request", ctx.Request.URL())

		// LoadResponse runs the default request to the destination of the request.
		// Not calling this will require you to mock the entire response.
		// This can be done with the SetXxx (Status, Header, Body) functions on the
		// ctx.Response struct.
		_ = ctx.LoadResponse(http.DefaultClient, true)

		uri := ctx.Request.URL().RequestURI()
		seccodes[uri] = ctx.Response.Body()

		// Here we append some code to every js file.
		// The code will update the document title to "hi"
		ctx.Response.SetBody(seccodes[uri])

		err := os.WriteFile(
			fmt.Sprintf("./captchas/%s.bmp", uri),
			[]byte(seccodes[uri]),
			0600,
		)
		if err != nil {
			slog.Error("write file err", "err", err)
		}
	})
	go router.Run()

	// Create a new page
	page := browser.MustPage("https://discuz34.dev.lan").MustWaitLoad()
	page.WaitStable(time.Millisecond * 800)

	// We use css selector to get the search input element and input "git"
	page.MustElement("#lsform #ls_username").MustInput("admin")
	page.MustElement("#lsform #ls_password").MustInput("123456")
	page.MustElement("#lsform button").Click(proto.InputMouseButtonLeft, 1)
	page.WaitStable(time.Millisecond * 2000)

	// Refresh captcha
	// e, err := page.Element("#fwin_login")
	// if err == nil {
	// 	e.MustElement(".fwin form a").Click(proto.InputMouseButtonLeft, 1)
	// 	page.WaitStable(time.Millisecond * 2000)
	// }

	slog.Info("Pause")
	utils.Pause()
}

func newBrowser() (*rod.Browser, func()) {
	l := launcher.New().
		Bin("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome").
		UserDataDir("tmp/user").
		Set("disable-default-apps").
		Set("no-first-run")
	if path, exists := launcher.LookPath(); exists && false {
		l = launcher.New().Bin(path)
	}
	l.Headless(false).Devtools(true)
	// defer l.Cleanup() // remove launcher.FlagUserDataDir

	// Launch a new browser with default options, and connect to it.
	// Trace shows verbose debug information for each action executed
	// Slowmotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	browser := rod.New().
		ControlURL(l.MustLaunch()).
		Trace(false).
		SlowMotion(2 * time.Second).
		MustConnect().
		NoDefaultDevice()

	// ServeMonitor plays screenshots of each tab. This feature is extremely
	// useful when debugging with headless mode.
	// You can also enable it with flag "-rod=monitor"
	launcher.Open(browser.ServeMonitor(""))

	// Even you forget to close, rod will close it after main process ends.
	// defer browser.MustClose()

	incognito, err := browser.Incognito()
	if err != nil {
		panic(err)
	}
	// defer incognito.Close()

	return incognito, func() {
		defer l.Cleanup() // remove launcher.FlagUserDataDir
		defer browser.MustClose()
		defer incognito.Close()
	}
}

// func GetCaptcha(reg, eSelector string, p *rod.Page) (captcha string) {
// 	var byteImg []byte
// 	r, _ := regexp.Compile(reg)

// 	for !r.MatchString(captcha) {
// 		p.MustElement(eSelector).MustClick().MustWaitStable()
// 		byteImg = p.MustElement(eSelector).MustResource()

// 		resp, err := http.Post(apiURL, "application/json;charset=UTF-8",
// 			strings.NewReader(base64.StdEncoding.EncodeToString(byteImg)))
// 		if err != nil {
// 			panic(err)
// 		}
// 		defer resp.Body.Close()
// 		data, _ := io.ReadAll(resp.Body)
// 		captcha = gjson.Parse(string(data)).Get("result").String()
// 	}

// 	log.Println("验证码获取成功:", captcha)
// 	return captcha
// }
