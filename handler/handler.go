package handler 

import (
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"net/url"
	"time"
	"urlShortener/storage"
)

func New(schema string, host string , storage storage.service) *router.Router{
	router := router.New()

	h := handler{schema, host , storage}
	router.POST("/encode", responseHandler(h.encode))
	router.GET("/{shortLink}", h.redirect)
	router.GET("/{shortLink}/info", responseHandler(h.decode))
	return  router;
}

type response struct {
	schema string
	host string
	storage  storage.service
}

func responseHandler(h func(ctx *fasthttp.RequestCtx) (interface{}, int, error)) fasthttp.RequestHandler {
	return func (ctx *fasthttp.RequestCtx) {
		data , status , err := h(ctx)
		if err != nil {
			data = err.Error()
		}
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.SetStatusCode(status)
		err = json.NewEncoder(ctx.Response.BodyWriter()).Encode(response{Data : data , Success : err == nill})
		if err != nil {
			log.Printf("could not encode response to output : %v", err)
		}
	}
}

func (h handler) encode (ctx *fasthttp.RequestCtx) (interface{})

