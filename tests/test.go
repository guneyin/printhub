package test

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/printhub/cmd/api"
	"io"
)

const testPort = "8081"

type testCase struct {
	skip               bool
	printOutput        bool
	description        string
	method             string
	route              string
	body               body
	expectedStatusCode int
}

type body struct {
	data any
}

func newBody(data any) body {
	return body{data: data}
}

func (b body) toBytes() []byte {
	d, _ := json.Marshal(b.data)
	return d
}

func (b body) toReader() io.Reader {
	return bytes.NewReader(b.toBytes())
}

var testApp *fiber.App

func init() {
	app, err := api.NewApplication()
	if err != nil {
		panic(err)
	}
	testApp = app.Server
}
