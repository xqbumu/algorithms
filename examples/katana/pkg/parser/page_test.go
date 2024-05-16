package parser_test

import (
	"algorithms/examples/katana/pkg/parser"
	"log"
	"testing"
)

func TestPageCdekExpress(t *testing.T) {
	g := setup(t)

	p := g.page.MustNavigate(g.srcFile("fixtures/cdek-express.html")) // .MustWaitLoad()

	rows, err := parser.CdekExpress(p)
	g.Nil(err)
	g.Len(rows, 9)

	log.Println(rows)
}

func TestPageSibtrans(t *testing.T) {
	g := setup(t)

	p := g.page.MustNavigate(g.srcFile("fixtures/sibtrans.html")) // .MustWaitLoad()
	rows, err := parser.CdekExpress(p)
	g.Nil(err)
	g.Len(rows, 5)

	log.Println(rows)
}
