package http

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
)

type APIResponseCode int
const (
	ERROR APIResponseCode = iota
	OK
)

type HttpAPI struct {
	URLPrefix string
}

var API HttpAPI = HttpAPI{}

type APIResponse struct {
	Code APIResponseCode
	Message string
}

func (this *APIResponseCode)HttpStatus() int {
	switch *this {
	case ERROR:
		return http.StatusInternalServerError
	case OK:
		return http.StatusOK
	}
	return http.StatusNotImplemented
}

func Respond(r render.Render, apiResponse *APIResponse)  {
	r.JSON(apiResponse.Code.HttpStatus(), apiResponse)
}

func (this *HttpAPI)registerAPIRequests(m *martini.ClassicMartini, path string, handler martini.Handler)  {
	fullpath := fmt.Sprintf("/api/%s", path)
	m.Get(fullpath, handler)
	//m.Post()
	//支持post以及get

}

func (this *HttpAPI)RegisterRequests(m *martini.ClassicMartini)  {
	log.Println("Enter RegisterRequests")
	//建立url和handler之间的关系
	this.registerAPIRequests(m, "test/:host/:port", this.GetHandler)
}

func (this *HttpAPI) GetHandler(params martini.Params, r render.Render, req *http.Request)  {

	Respond(r, &APIResponse{Code:OK, Message:fmt.Sprintf("Test:%s %s",params["host"], params["port"])})
}
