package http

import (
	"Hermes/internal/providers"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
)

func (i *HTTPInstanceAPI) addOrder(ctx *fasthttp.RequestCtx) {
	var order providers.Order
	body := ctx.Request.Body()
	if err := json.Unmarshal(body, &order); err != nil {
		i.log.Errorf("Unable to unmarshal response: %v", err)
		ctx.Response.SetBodyString("Invalid request sent")
		ctx.Response.SetStatusCode(400)
		return
	}

	if err := i.api.AddOrder(&order); err != nil {
		i.log.Errorf("Unable to add order: %v", err)
		ctx.Response.SetBodyString("Invalid request sent")
		ctx.Response.SetStatusCode(400)
		return
	}
	handleRequest(ctx)
	ctx.Response.SetStatusCode(201)
}

func (i *HTTPInstanceAPI) getOrders(ctx *fasthttp.RequestCtx) {
	result, err := i.api.GetOrders(ctx)
	if err != nil {
		ctx.Response.SetBodyString("Cannot getting orders")
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

func (i *HTTPInstanceAPI) getOrder(ctx *fasthttp.RequestCtx) {
	userOrderID := ctx.UserValue("order_id").(string)
	orderID, err := strconv.ParseInt(userOrderID, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("error while parsing data: %v", err)
		ctx.Response.SetBodyString(msg)
		ctx.Response.SetStatusCode(400)
		return
	}

	order, err := i.api.GetOrder(orderID)
	if err != nil {
		msg := fmt.Sprintf("error while retriving from mongo: %v", err)
		ctx.Response.SetBodyString(msg)
		ctx.Response.SetStatusCode(400)
		return
	}

	body, err := json.Marshal(order)
	if err != nil {
		msg := fmt.Sprintf("error while retriving from mongo: %v", err)
		ctx.Response.SetBodyString(msg)
		ctx.Response.SetStatusCode(400)
		return
	}
	ctx.Response.SetBody(body)
	ctx.Response.SetStatusCode(200)

}

func (i *HTTPInstanceAPI) deleteOrder(ctx *fasthttp.RequestCtx) {
	userOrderID := ctx.UserValue("order_id").(string)
	orderID, err := strconv.ParseInt(userOrderID, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("error while parsing data: %v", err)
		ctx.Response.SetBodyString(msg)
		ctx.Response.SetStatusCode(400)
		return
	}
	if err := i.api.DeleteOrder(orderID); err != nil {
		msg := fmt.Sprintf("error while deleting order: %v", err)
		ctx.Response.SetBodyString(msg)
		ctx.Response.SetStatusCode(400)
		return
	}
	ctx.Response.SetStatusCode(200)
	handleRequest(ctx)
}
