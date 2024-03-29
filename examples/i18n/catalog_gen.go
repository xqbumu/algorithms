// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p, ok := messageKeyToIndex[key]
	if !ok {
		return "", false
	}
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"en_US": &dictionary{index: en_USIndex, data: en_USData},
		"zh":    &dictionary{index: zhIndex, data: zhData},
	}
	fallback := language.MustParse("en-US")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"Welcome %s!": 0,
}

var en_USIndex = []uint32{ // 2 elements
	0x00000000, 0x0000000f,
} // Size: 32 bytes

const en_USData string = "\x02Welcome %[1]s!"

var zhIndex = []uint32{ // 2 elements
	0x00000000, 0x00000000,
} // Size: 32 bytes

const zhData string = ""

// Total table size 79 bytes (0KiB); checksum: 72BAF01C
