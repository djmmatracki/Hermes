package internal

import (
	"encoding/json"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"gopkg.in/validator.v2"
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

	i.log.Infof("Starting server at port %s", i.bind)
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
	var body SingleLaunchRequest
	collection := i.api.mongoDatabase.Collection("trucks")

	err := json.Unmarshal(ctx.Request.Body(), &body)

	if err != nil {
		i.log.Infof("Unable to unmarshal response: %v", err)
		ctx.Response.SetBodyString("Invalid request sent")
		ctx.Response.SetStatusCode(400)
	}

	// Data validation
	if err := validator.Validate(body); err != nil {
		ctx.Response.SetBodyString("Invelid input data")
		ctx.Response.SetStatusCode(400)
	} else {

		// Execute insertion
		_, err := collection.InsertOne(ctx, body)
		if err != nil {
			ctx.Response.SetBodyString("Error while inserting truck")
			ctx.Response.SetStatusCode(400)
		} else {
			ctx.Response.SetBodyString("Inserted new truck...")
			ctx.Response.SetStatusCode(200)
		}
	}
}

func (i *HTTPInstanceAPI) getTrucks(ctx *fasthttp.RequestCtx) {
	result, err := i.api.getTrucks(ctx)
	if err != nil {
		ctx.Response.SetBodyString("Cannot get trucks")
		ctx.Response.SetStatusCode(400)
		return
	}
	body, _ := json.Marshal(result)
	ctx.Response.SetBody(body)
	ctx.Response.SetStatusCode(200)
}

func (i *HTTPInstanceAPI) handleRoot(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBodyString("Welcome to Hermes!")
}
