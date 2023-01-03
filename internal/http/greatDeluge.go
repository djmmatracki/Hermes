package http

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type MatchOrdersRequestDeluge struct {
	WaterLevel  float64 `json:"water_level"`
	RainSpeed   float64 `json:"rain_speed"`
	GroundLevel float64 `json:"ground_level"`
	NumberLoops float64 `json:"number_loops"`
}

func (i *HTTPInstanceAPI) greatDeluge(ctx *fasthttp.RequestCtx) {
	var matchOrdersRequest MatchOrdersRequestDeluge
	body := ctx.Request.Body()
	if err := json.Unmarshal(body, &matchOrdersRequest); err != nil {
		i.log.Errorf("could't parse body in request")
		return
	}

	solution, err := i.api.GreatDelugeAlgorithm(
		matchOrdersRequest.WaterLevel,
		matchOrdersRequest.RainSpeed,
		matchOrdersRequest.GroundLevel,
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
