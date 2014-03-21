Goztivo
=======

A web API wrapper for [Oztivo](http://www.oztivo.net) - a TV guide database for Australia

Goztivo has 2 components:

* Javascript/HTML frontend which uses AngularJS framework
* an API server backend written in Golang. This also doubles as a general purpose web server
to serve the HTML/JS content

Goztivo aims to follow Oztivo's usage policy where possible by sending a unique HTTP 
user-agent string and keeping traffic to a minimum by caching data.

## Quick start

To run Goztivo, first [install Golang](http://golang.org/doc/install#install). On Redhat Linux that should be as simple as `yum install golang`. Then:

- Install the Goztivo repository `go get github.com/porjo/goztivo`
- Change to Goztivo's install loction `cd $GOPATH/src/github.com/porjo/goztivo`
- Run Goztivo `go run main.go`
- Browse to your IP on port `3000` e.g. `http://localhost:3000`

## Cache

Goztivo caches data on disk in the default temporary directory location (e.g. `/tmp`) with directory prefix `goztivo_`. It will store upto 100MB of data by default.

## API

The API is very simple consisting of 2 endpoints: `/api/programme` and `/api/channel`. Each accepts a POST request containing JSON body and returns a JSON encoded response.

**Channels**

To retrieve all available channels, simply send an empty request:
```
curl --data '' http://localhost:3000/api/channel
```

To retrieve all channels containing the substring *qld* (case insensitive):
```
curl -i --data '{"channel_name": "qld", "contains": true}' http://localhost:3000/api/channel
```

**Programmes**

To retrieve all programmes for a given channel on a given day, supply an array of `channel_id`s (retrieved previously) and an array of `day`s in JSON timestamp format:
```
curl -i --data '{"channels": ["WIN-Qld"], "days": ["2014-03-13T00:25:43.511Z"]}' http://localhost:3000/api/programme
```

## Licence

Goztivo is under The MIT license. TV guide data from Oztivo is under Creative Commons ([CC by-nc-sa](http://creativecommons.org/licenses/by-nc-sa/2.5/au/)) and not for commercial use.

