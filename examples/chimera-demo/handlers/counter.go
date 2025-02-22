package handlers

import (
	"fmt"

	"github.com/matt1484/chimera"
)

type CounterBody struct {
	Count int    `json:"count"`
	Msg   string `json:"message"`
}

type CounterParams struct {
	Name string `param:"name,in=query"`
}

type Counter struct {
	count int
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) DoHandler(req *chimera.FormRequest[CounterBody, CounterParams]) (*chimera.JSON[CounterBody, CounterParams], error) {
	c.count++
	body := CounterBody{
		Count: c.count,
	}
	if req.Params.Name != "" {
		body.Msg = fmt.Sprintf("Hello, %s!", req.Params.Name)
	}
	return &chimera.JSON[CounterBody, CounterParams]{
		Body: body,
	}, nil
}
