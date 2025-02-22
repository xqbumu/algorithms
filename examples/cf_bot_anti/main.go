package main

import (
	"algorithms/examples/scrapify/pkg"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func main() {
	var err error

	// Setup your client however you need it. This is simply an example
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   15 * time.Second,
				KeepAlive: 15 * time.Second,
				DualStack: true,
			}).DialContext,
		},
	}
	// Set the client Transport to the RoundTripper that solves the Cloudflare anti-bot
	client.Transport, err = pkg.NewCFRT(client.Transport)
	if err != nil {
		return
	}

	req, err := http.NewRequest("GET", "https://www.sgpbusiness.com", nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

// func main() {
// 	u := os.Args[1]
// 	cl := pkg.TlsRequest()
// 	req, err := http.NewRequest("GET", u, nil)

// 	// user agent must be set
// 	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`)
// 	resp, err := cl.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(body))
// }

// func main() {
// 	client := &http.Client{
// 		Transport: scrapify.NewTransport(scrapify.Firefox),
// 	}
// 	req, err := http.NewRequest(http.MethodGet, "https://www.sgpbusiness.com", nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	scrapify.SetHeaders(req, nil)
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println(resp.StatusCode)
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(body))
// }
