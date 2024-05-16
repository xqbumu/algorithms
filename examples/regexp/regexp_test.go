package regexp_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexpMultiLine(t *testing.T) {
	var re = regexp.MustCompile(`(?sm)<p.*>(.*)<\/p>`)
	var str = `<div><p class='block_paragraph'>
Заказ завершен
06.08.2023.

</p>`

	result := re.FindStringSubmatch(str)
	assert.Equal(t, 2, len(result))

	fmt.Println(result[1]) // final resule
}
