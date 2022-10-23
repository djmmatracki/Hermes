package internal

import (
	"fmt"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type HTTPInstanceAPI struct {
	bind string
	log  logrus.FieldLogger
}

func NewHTTPInstanceAPI(bind string) *HTTPInstanceAPI {
	return &HTTPInstanceAPI{
		bind: bind,
	}
}

func (i *HTTPInstanceAPI) Run() {
	r := router.New()

	r.GET("/", i.handleRoot)

	fasthttpServer := fasthttp.Server{}
	err := fasthttpServer.ListenAndServe(i.bind)
	fmt.Printf("Error occured %v", err)
}

func (i *HTTPInstanceAPI) handleRoot(ctx *fasthttp.RequestCtx) {
	ListPoints()
	ctx.Response.SetBodyString("Welcome to Hermes!")
}
