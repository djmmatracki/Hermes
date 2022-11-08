package internal

import (
	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type HTTPInstanceAPI struct {
	bind string
	log  logrus.FieldLogger
	api  *InstanceAPI
}

func NewHTTPInstanceAPI(bind string, log logrus.FieldLogger, api *InstanceAPI) *HTTPInstanceAPI {
	return &HTTPInstanceAPI{
		bind: bind,
		log:  log,
		api:  api,
	}
}

func (i *HTTPInstanceAPI) Run() {
	r := router.New()

	// Root endpoint
	r.GET("/", i.handleRoot)
	r.GET("/a-star", i.aStar)

	// Truck endpoints
	r.POST("/truck", i.addTruck)
	r.GET("/truck", i.getTrucks)

	// Generate optimal
	r.POST("/launch", i.launch)

	i.log.Fatal(fasthttp.ListenAndServe(i.bind, r.Handler))
}

func (i *HTTPInstanceAPI) launch(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBodyString("Genereting optimal paths for given fleet...")
	// Receive list of orders
	// Loop each truck and each order
	// For each pair (truck - order) calculate priority
}
func (i *HTTPInstanceAPI) aStar(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBodyString("Compiling a-star...")
}

func (i *HTTPInstanceAPI) addTruck(ctx *fasthttp.RequestCtx) {
	// Get response from api
	// Parse response
	// Execute insertion
	ctx.Response.SetBodyString("Inserted new truck...")
}

func (i *HTTPInstanceAPI) getTrucks(ctx *fasthttp.RequestCtx) {
	i.api.getTrucks(ctx)
	ctx.Response.SetBodyString("Get all truck...")
}

func (i *HTTPInstanceAPI) handleRoot(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBodyString("Welcome to Hermes!")
}
