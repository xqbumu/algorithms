package handlers

import "github.com/matt1484/chimera"

type TestBody struct {
	Property string `json:"prop"`
}

type TestParams struct {
	Path   string `param:"path,in=path"`
	Header string `param:"header,in=header"`
}

func Test(req *chimera.JSON[TestBody, TestParams]) (*chimera.JSON[TestBody, TestParams], error) {
	return &chimera.JSON[TestBody, TestParams]{
		Body: req.Body,
	}, nil
}
