// goztivo provides a grabber for Oztivo data
//
// HTTP client caches results as per Oztivo API usage policy:
// http://www.oztivo.net/twiki/bin/view/TVGuide/StaticXMLGuideAPI
//
package main

import (
	"log"

	"github.com/codegangsta/martini"
	"github.com/porjo/goztivo/api"
)

const (
	NAME    = "Goztivo"
	VERSION = "0.1"
)

func main() {
	log.Println(NAME + " Starting")

	err := api.InitAPI(NAME + "/" + VERSION)
	if err != nil {
		log.Fatal(err)
	}
	m := martini.Classic()
	m.Post("/api/programme", api.ProgrammeHandler)
	m.Post("/api/channel", api.ChannelHandler)
	m.Run()
}
