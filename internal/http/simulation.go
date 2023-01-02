package http

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

// func (i *HTTPInstanceAPI) singleLaunch(ctx *fasthttp.RequestCtx) {
// 	var singleLaunchRequst SingleLaunchRequest
// 	body := ctx.Request.Body()

// 	err := json.Unmarshal(body, &singleLaunchRequst)
// 	if err != nil {
// 		i.log.Infof("Unable to unmarshal response: %v", err)
// 		ctx.Response.SetBodyString("Invalid request sent")
// 		ctx.Response.SetStatusCode(400)
// 		return
// 	}

// 	response, err := i.api.SingleTruckLaunch(
// 		singleLaunchRequst.TruckID,
// 		Location{
// 			Latitude:  singleLaunchRequst.OriginLat,
// 			Longitude: singleLaunchRequst.OriginLon,
// 		},
// 		Location{
// 			Latitude:  singleLaunchRequst.DestinationLat,
// 			Longitude: singleLaunchRequst.DestinationLon,
// 		},
// 		Location{
// 			Latitude:  singleLaunchRequst.DestinationLat,
// 			Longitude: singleLaunchRequst.DestinationLon,
// 		},
// 	)
// 	if err != nil {
// 		i.log.Infof("Unable to unmarshal response: %v", err)
// 		ctx.Response.SetBodyString("Invalid request sent")
// 		ctx.Response.SetStatusCode(400)
// 		return
// 	}

// 	jsonResponse, err := json.Marshal(response)
// 	if err != nil {
// 		i.log.Infof("Unable to unmarshal response: %v", err)
// 		ctx.Response.SetBodyString("Invalid request sent")
// 		ctx.Response.SetStatusCode(500)
// 	}
// 	ctx.Response.SetBody(jsonResponse)
// 	ctx.Response.SetStatusCode(200)
// }

type MatchOrdersRequest struct {
	MaxIterations     int     `json:"max_iterations"`
	StartTemperature  float64 `json:"start_temperature"`
	FinishTemperature float64 `json:"finish_temperature"`
	Cooling           float64 `json:"cooling"`
}

func (i *HTTPInstanceAPI) handleRoot(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBodyString("Welcome to Hermes!")
}

func (i *HTTPInstanceAPI) simulate(ctx *fasthttp.RequestCtx) {
	var matchOrdersRequest MatchOrdersRequest
	body := ctx.Request.Body()
	if err := json.Unmarshal(body, &matchOrdersRequest); err != nil {
		i.log.Errorf("could't parse body in request")
		return
	}

	solution, err := i.api.SimulatedAnneling(
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
