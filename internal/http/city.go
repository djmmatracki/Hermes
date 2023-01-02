package http

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

func (i *HTTPInstanceAPI) getCities(ctx *fasthttp.RequestCtx) {
	result, err := i.api.GetCities(ctx)
	if err != nil {
		ctx.Response.SetBodyString("Cannot get cities")
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
