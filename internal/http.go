package internal

import (
	"encoding/json"

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
	r.POST("/single-launch", i.singleLaunch)

	i.log.Fatal(fasthttp.ListenAndServe(i.bind, r.Handler))
}

func (i *HTTPInstanceAPI) singleLaunch(ctx *fasthttp.RequestCtx) {
	var singleLaunchRequst SingleLaunchRequest

	body := ctx.Request.Body()
	err := json.Unmarshal(body, &singleLaunchRequst)

	if err != nil {
		i.log.Infof("Unable to unmarshal response: %v", err)
		ctx.Response.SetBodyString("Invalid request sent")
		ctx.Response.SetStatusCode(400)
		return
	}
	response, err := i.api.singleTruckLaunch(
		singleLaunchRequst.TruckID,
		Location{
			Latitude:  singleLaunchRequst.OriginLat,
			Longitude: singleLaunchRequst.OriginLon,
		},
		Location{
			Latitude:  singleLaunchRequst.DestinationLat,
			Longitude: singleLaunchRequst.DestinationLon,
		},
	)
	if err != nil {
		i.log.Infof("Unable to unmarshal response: %v", err)
		ctx.Response.SetBodyString("Invalid request sent")
		ctx.Response.SetStatusCode(400)
		return
	}
	jsonResponse, err := json.Marshal(response)

	if err != nil {
		i.log.Infof("Unable to unmarshal response: %v", err)
		ctx.Response.SetBodyString("Invalid request sent")
		ctx.Response.SetStatusCode(500)
	}
	ctx.Response.SetBody(jsonResponse)
	ctx.Response.SetStatusCode(200)
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
