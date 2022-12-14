package internal

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"gopkg.in/validator.v2"
)

type MatchOrdersRequest struct {
	Orders            []Order `json:"orders"`
	MaxIterations     int     `json:"max_iterations"`
	StartTemperature  float32 `json:"start_temperature"`
	FinishTemperature float32 `json:"finish_temperature"`
	Cooling           float32 `json:"cooling"`
}

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

	r.GET("/", i.handleRoot)

	r.POST("/truck", i.addTruck)
	r.GET("/truck", i.getTrucks)
	r.GET("/truck/{truck_id}", i.getTruck)

	r.POST("/single-launch", i.singleLaunch)
	r.POST("/match-orders", i.matchOrders)
	r.POST("/a-star", i.Astar)

	i.log.Infof("Starting server at port %s", i.bind)
	i.log.Fatal(fasthttp.ListenAndServe(i.bind, r.Handler))
}

func (i *HTTPInstanceAPI) Astar(ctx *fasthttp.RequestCtx) {
	mainCollection := i.api.mongoDatabase.Collection("main")
	dist, err := Astar(
		mainCollection,
		Location{
			Latitude:  0.2,
			Longitude: 0.2,
		},
		Location{
			Latitude:  0.2,
			Longitude: 0.2,
		},
	)
	if err != nil {
		return
	}

	i.log.Infof("distance: %f", dist)
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

func (i *HTTPInstanceAPI) addTruck(ctx *fasthttp.RequestCtx) {

	// Get response from api
	var newTruck Truck
	collection := i.api.mongoDatabase.Collection("truck")

	body := ctx.Request.Body()
	err := json.Unmarshal(body, &newTruck)

	if err != nil {
		i.log.Infof("Unable to unmarshal response: %v", err)
		ctx.Response.SetBodyString("Invalid request sent")
		ctx.Response.SetStatusCode(400)
		return
	}

	// Data validation
	err2 := validator.Validate(newTruck)
	if err2 != nil {
		ctx.Response.SetBodyString("Invelid input data")
		ctx.Response.SetStatusCode(400)
		return
	}

	// Execute insertion
	_, err3 := collection.InsertOne(ctx, newTruck)
	if err3 != nil {
		ctx.Response.SetBodyString("Error while inserting truck")
		ctx.Response.SetStatusCode(400)
		return
	} else {
		ctx.Response.SetBodyString("Inserted new truck...")
		ctx.Response.SetStatusCode(200)
		return
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

func (i *HTTPInstanceAPI) matchOrders(ctx *fasthttp.RequestCtx) {
	var matchOrdersRequest MatchOrdersRequest
	body := ctx.Request.Body()
	if err := json.Unmarshal(body, &matchOrdersRequest); err != nil {
		i.log.Errorf("could't parse body in request")
		return
	}

	solution, err := i.api.simulatedAnneling(
		matchOrdersRequest.Orders,
		matchOrdersRequest.MaxIterations,
		matchOrdersRequest.StartTemperature,
		matchOrdersRequest.FinishTemperature,
		matchOrdersRequest.Cooling,
		10,
	)
	if err != nil {
		i.log.Errorf("error occured while executing anneling algorithm %v", err)
		return
	}

	responseBody, err := json.Marshal(solution)
	if err != nil {
		i.log.Errorf("error occured while marshaling solution %v", err)
		return
	}
	ctx.Response.SetBody(responseBody)
	ctx.Response.SetStatusCode(200)
}

func (i *HTTPInstanceAPI) getTruck(ctx *fasthttp.RequestCtx) {
	userTruckID := ctx.UserValue("truck_id").(string)
	truckID, err := strconv.ParseInt(userTruckID, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("error while parsing data: %v", err)
		ctx.Response.SetBodyString(msg)
		ctx.Response.SetStatusCode(400)
		return
	}

	truck, err := i.api.getTruck(ctx, TruckID(truckID))
	if err != nil {
		msg := fmt.Sprintf("error while retriving from mongo: %v", err)
		ctx.Response.SetBodyString(msg)
		ctx.Response.SetStatusCode(400)
		return
	}
	body, err := json.Marshal(truck)
	if err != nil {
		msg := fmt.Sprintf("error while retriving from mongo: %v", err)
		ctx.Response.SetBodyString(msg)
		ctx.Response.SetStatusCode(400)
		return
	}
	ctx.Response.SetBody(body)
	ctx.Response.SetStatusCode(200)
}
