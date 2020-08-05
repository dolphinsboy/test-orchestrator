package main

import (
	"github.com/dolphinsboy/test-orchestrator/go/http"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	nethttp "net/http"
)

func main() {
	m := martini.Classic()
	http.API.URLPrefix = ""
	http.Web.URLPrefix = ""
	http.API.RegisterRequests(m)
	http.Web.RegisterRequests(m)

	//fixed https://github.com/martini-contrib/render/issues/43
	m.Use(render.Renderer(render.Options{
		Directory:"resources",
		HTMLContentType:"text/html",
	}))

	if err := nethttp.ListenAndServe(":8000", m); err != nil {
		log.Fatal(err)
	}
}

