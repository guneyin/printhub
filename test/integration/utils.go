package integration

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/cmd/api"
)

type Case struct {
	Skip               bool
	PrintOutput        bool
	Description        string
	Method             string
	Route              string
	Body               Body
	ExpectedStatusCode int
}

type Body struct {
	data any
}

func NewBody(data any) Body {
	return Body{data: data}
}

func (b Body) toBytes() []byte {
	d, _ := json.Marshal(b.data)
	return d
}

func (b Body) ToReader() io.Reader {
	return bytes.NewReader(b.toBytes())
}

func NewTestApp() *fiber.App {
	app, err := api.NewApplication()
	if err != nil {
		panic(err)
	}

	return app.Server
}
