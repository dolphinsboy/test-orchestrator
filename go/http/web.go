package http

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
)

type Result struct {
	Title string
	Content string
}

type HttpWeb struct {
	URLPrefix string
}

var Web HttpWeb = HttpWeb{
}

func (this *HttpWeb)registerWebRequest(m *martini.ClassicMartini, path string, handler martini.Handler)  {
	fullpath := fmt.Sprintf("%s/web/%s", this.URLPrefix, path)
	log.Println(fullpath)
	m.Get(fullpath, handler)
}

func (this *HttpWeb)Test(params martini.Params, r render.Render, req *http.Request)  {
	ret := Result{
		Title:"test",
		Content:"Values",
	}
	r.HTML(200, "templates/test",ret)
}

func (this *HttpWeb)RegisterRequests(m *martini.ClassicMartini)  {
	this.registerWebRequest(m, "test", this.Test)
}
