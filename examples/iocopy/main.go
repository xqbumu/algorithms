package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
)

func main() {
	body1 := bytes.NewBuffer(nil)
	fmt.Fprintf(body1, "hello")

	body2 := bytes.NewBuffer(nil)
	reader := io.TeeReader(body1, body2)

	data, err := io.ReadAll(reader)
	log.Printf("body1 %s %v", data, err)

	data2, err := io.ReadAll(body2)
	log.Printf("body2 %s %v", data2, err)

	s := "ewogICAgImVycm9yIjogewogICAgICAgICJtZXNzYWdlIjogIllvdSBleGNlZWRlZCB5b3VyIGN1cnJlbnQgcXVvdGEsIHBsZWFzZSBjaGVjayB5b3VyIHBsYW4gYW5kIGJpbGxpbmcgZGV0YWlscy4iLAogICAgICAgICJ0eXBlIjogImluc3VmZmljaWVudF9xdW90YSIsCiAgICAgICAgInBhcmFtIjogbnVsbCwKICAgICAgICAiY29kZSI6IG51bGwKICAgIH0KfQo="
	body, err := base64.StdEncoding.DecodeString(s)
	log.Printf("%s", body)
}
