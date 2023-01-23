package http

import (
	"Hermes/internal/providers"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
)

func (i *HTTPInstanceAPI) getTruck(ctx *fasthttp.RequestCtx) {
	userTruckID := ctx.UserValue("truck_id").(string)
	truckID, err := strconv.ParseInt(userTruckID, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("error while parsing data: %v", err)
		ctx.Response.SetBodyString(msg)
		ctx.Response.SetStatusCode(400)
		return
	}
	// Added comment

	truck, err := i.api.GetTruck(ctx, truckID)
	if err != nil {
		msg := fmt.Sprintf("error while retriving from mongo: %v", err)
		ctx.Response.SetBodyString(msg)
		ctx.Response.SetStatusCode(400)
		return
	}
	// Added different comment

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

func (i *HTTPInstanceAPI) getTrucks(ctx *fasthttp.RequestCtx) {
	result, err := i.api.GetTrucks(ctx)
	if err != nil {
		ctx.Response.SetBodyString("Cannot get trucks")
		ctx.Response.SetStatusCode(500)
		return
	}
	body, err := json.Marshal(result)
	if err != nil {
		ctx.Response.SetBodyString("Error while marshaling")
		ctx.Response.SetStatusCode(500)
		return
	}
	ctx.Response.SetBody(body)
	ctx.Response.SetStatusCode(200)
	handleRequest(ctx)
}

func (i *HTTPInstanceAPI) addTruck(ctx *fasthttp.RequestCtx) {
	var truck providers.Truck
	body := ctx.Request.Body()
	if err := json.Unmarshal(body, &truck); err != nil {
		i.log.Infof("Unable to unmarshal response: %v", err)
		ctx.Response.SetBodyString("Invalid request sent")
		ctx.Response.SetStatusCode(400)
		return
	}

	if err := i.api.AddTruck(truck); err != nil {
		i.log.Errorf("Got error while adding truck %v", err)
		ctx.Response.SetBodyString("Error while inserting truck")
		ctx.Response.SetStatusCode(400)
		return
	}

	ctx.Response.SetBodyString("Success")
	ctx.Response.SetStatusCode(200)
	handleRequest(ctx)
}

func (i *HTTPInstanceAPI) deleteTruck(ctx *fasthttp.RequestCtx) {
	userTruckID := ctx.UserValue("truck_id").(string)
	truckID, err := strconv.ParseInt(userTruckID, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("error while parsing data: %v", err)
		ctx.Response.SetBodyString(msg)
		ctx.Response.SetStatusCode(400)
		return
	}
	if err := i.api.DeleteTruck(truckID); err != nil {
		msg := fmt.Sprintf("error while deleting order: %v", err)
		ctx.Response.SetBodyString(msg)
		ctx.Response.SetStatusCode(400)
		return
	}
	ctx.Response.SetStatusCode(200)
	handleRequest(ctx)
}
