package internal

import (
	"log"

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

	// Root endpoint
	r.GET("/", i.handleRoot)

	// Truck endpoints
	r.POST("truck/", i.addTruck)
	r.GET("truck/", i.getTruck)

	// Generate optimal
	r.POST("launch/", i.launch)

	log.Fatal(fasthttp.ListenAndServe(i.bind, r.Handler))
}

func (i *HTTPInstanceAPI) launch(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBodyString("Genereting optimal paths for given fleet...")
	// Receive list of orders
	// Loop each truck and each order
	// For each pair (truck - order) calculate priority
}

func (i *HTTPInstanceAPI) addTruck(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBodyString("Welcome to Hermes!")
}

func (i *HTTPInstanceAPI) getTruck(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBodyString("Welcome to Hermes!")
}

func (i *HTTPInstanceAPI) calcuteSingle(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBodyString("Welcome to Hermes!")
}

func (i *HTTPInstanceAPI) handleRoot(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBodyString("Welcome to Hermes!")
}
