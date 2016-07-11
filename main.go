// Goztivo: a web API wrapper for Oztivo, a TV guide database for Australia
//
// Goztivo has 2 components:
//
// - Javascript/HTML frontend which uses AngularJS framework
// - an API server backend written in Golang. This also doubles as a general
// purpose web server to serve the HTML/JS content
//
// Goztivo aims to follow Oztivo's usage policy where possible by sending a
// unique HTTP user-agent string and keeping traffic to a minimum by caching
// data and limiting upstream requests to a maximum of one HTTP connection per second.

package main

import (
	"log"

	"github.com/codegangsta/martini"
	"github.com/porjo/goztivo/api"
)

const (
	NAME    = "Goztivo"
	VERSION = "0.3"
)

func main() {
	log.Println(NAME + " Starting")

	agentString := NAME + "/" + VERSION
	api.InitAPI(agentString)

	m := martini.Classic()
	m.Post("/api/programme", api.ProgrammeHandler)
	m.Post("/api/channel", api.ChannelHandler)
	m.Run()
}
