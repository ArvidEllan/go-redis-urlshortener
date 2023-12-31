package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"GOREDIS/storage"
)

func New(schema string, host string, storage storage.service) *router.Router {
	router := router.New()

	h := Handler{schema, host, storage}
	router.POST("/encode", responseHandler(h.encode))
	router.GET("/{shortLink}", h.redirect)
	router.GET("/{shortLink}/info", responseHandler(h.decode))
	return router
}

type response struct {
	schema  string
	host    string
	storage storage.service
}

func responseHandler(h func(ctx *fasthttp.RequestCtx) (interface{}, int, error)) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		data, status, err := h(ctx)
		if err != nil {
			data = err.Error()
		}
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.SetStatusCode(status)
		err = json.NewEncoder(ctx.Response.BodyWriter()).Encode(response{data})
		if err != nil {
			log.Printf("could not encode response to output : %v", err)
		}
	}
}

func (h handler) encode(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	var input struct {
		URL     string `json:"url"`
		Expires string `json:"expires"`
	}
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to decode json request body :%v", err)
	}

	uri, err := url.ParseRequestURI(input.URL)

	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid URL")
	}

	layoutISO := "2006-01-02 15:04:05"
	expires, err := time.Parse(layoutISO, input.Expires)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid expiration date")
	}

	c, err := h.storage.Save(uri.String(), expires)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("could not store in database: %v", err)
	}

	u := url.URL{
		Scheme: h.schema,
		Host:   h.host,
		Path:   c}
	fmt.Printf("Generated Link: %v \n", u.String())

	return u.String(), http.StatusCreated, nil
}
func (h handler) decode(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	code := ctx.UserValue("shortlink").(string)

	model, err := h.storage.LoadInfo(code)
	if err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.SetStatusCode(http.StatusNotFound)
		return
	}
	ctx.Redirect(uri, http.StatusMoved)
}
