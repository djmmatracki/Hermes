package http

import (
	"Hermes/internal/api"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type SingleLaunchRequest struct {
	TruckID        int     `json:"truck_id"`
	OriginLat      float32 `json:"origin_lat"`
	OriginLon      float32 `json:"origin_lon"`
	DestinationLat float32 `json:"destination_lat"`
	DestinationLon float32 `json:"destination_lon"`
}

type HTTPInstanceAPI struct {
	bind string
	log  logrus.FieldLogger
	api  *api.InstanceAPI
}

func NewHTTPInstanceAPI(bind string, log logrus.FieldLogger, api *api.InstanceAPI) *HTTPInstanceAPI {
	return &HTTPInstanceAPI{
		bind: bind,
		log:  log,
		api:  api,
	}
}

func (i *HTTPInstanceAPI) OptionsHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Add("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Response.Header.Add("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE")
	ctx.Response.Header.Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	ctx.Response.Header.Add("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.Add("Access-Control-Expose-Headers", "*")
	ctx.Response.SetStatusCode(200)
}

func (i *HTTPInstanceAPI) Run() {
	r := router.New()
	api := r.Group("/api")

	api.GET("/", i.handleRoot)
	api.OPTIONS("/{any}", i.OptionsHandler)

	api.POST("/truck", i.addTruck)
	api.GET("/truck", i.getTrucks)
	api.GET("/truck/{truck_id}", i.getTruck)
	api.DELETE("/truck/{truck_id}", i.deleteTruck)
	api.OPTIONS("/truck/{truck_id}", i.OptionsHandler)

	api.POST("/order", i.addOrder)
	api.DELETE("/order/{order_id}", i.deleteOrder)
	api.OPTIONS("/order/{order_id}", i.OptionsHandler)
	api.GET("/order", i.getOrders)
	api.GET("/order/{order_id}", i.getOrder)

	api.GET("/city", i.getCities)

	api.POST("/simulate", i.simulate)
	api.POST("/threshold", i.threshold)

	i.log.Infof("Starting server at port %s", i.bind)
	i.log.Fatal(fasthttp.ListenAndServe(i.bind, r.Handler))
}

func handleRequest(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Add("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Add("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE")
	ctx.Response.Header.Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	ctx.Response.Header.Add("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.Add("Access-Control-Expose-Headers", "*")
	ctx.Response.Header.Add("Content-type", "application/json charset=utf-8")
}
