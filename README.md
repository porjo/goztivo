goztivo
=======

A web API wrapper for [Oztivo](http://www.oztivo.net) - a TV guide database for Australia

Goztivo has 2 components:

* Javascript/HTML frontend which uses AngularJS framework
* an API server backend written in Golang. This also doubles as a general purpose web server
to server the HTML/JS content

Goztivo aims to follow Oztivo's usage policy where possible by sending a unique HTTP user-agent string
and keeping traffic to a minimum by caching data

## Quick start

To run Goztivo on Linux, first install Golang. On Redhat systems that should be as simple as `yum install golang`

- Install the Goztivo repository `go get http://github.com/porjo/goztivo`
- Change to goztivo loction `cd $GOPATH/src/goztivo`
- Run Goztivo `go run main.go`
- Browse to your IP on port `3000` e.g. `http://localhost:3000`


**NOTE:** This project is alpha 
