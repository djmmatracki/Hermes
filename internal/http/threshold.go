package http

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type MatchOrdersRequestThreshold struct {
	Threshold    int `json:"threshold"`
	ThrReduction int `json:"thr_reduction"`
	NumberLoops  int `json:"number_loops"`
}

func (i *HTTPInstanceAPI) threshold(ctx *fasthttp.RequestCtx) {
	var matchOrdersRequest MatchOrdersRequestThreshold
	body := ctx.Request.Body()
	if err := json.Unmarshal(body, &matchOrdersRequest); err != nil {
		i.log.Errorf("could't parse body in request")
		return
	}

	solution, err := i.api.ThresholdAccepting(
		matchOrdersRequest.Threshold,
		matchOrdersRequest.ThrReduction,
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
