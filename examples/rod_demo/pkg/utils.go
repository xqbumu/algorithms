package pkg

import (
	"fmt"
	"strings"

	"github.com/go-rod/rod"
)

func PrintReport(page *rod.Page) {
	el := page.MustElement("#broken-image-dimensions.passed")
	for _, row := range el.MustParents("table").First().MustElements("tr:nth-child(n+2)") {
		cells := row.MustElements("td")
		key := cells[0].MustProperty("textContent")
		if strings.HasPrefix(key.String(), "User Agent") {
			fmt.Printf("\t\t%s: %t\n\n", key, !strings.Contains(cells[1].MustProperty("textContent").String(), "HeadlessChrome/"))
		} else if strings.HasPrefix(key.String(), "Hairline Feature") {
			// Detects support for hidpi/retina hairlines, which are CSS borders with less than 1px in width, for being physically 1px on hidpi screens.
			// Not all the machine suppports it.
			continue
		} else {
			fmt.Printf("\t\t%s: %s\n\n", key, cells[1].MustProperty("textContent"))
		}
	}

	page.MustScreenshot("")
}
