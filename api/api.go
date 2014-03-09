package api

import (
	"encoding/json"
	//	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"code.google.com/p/go.text/encoding/charmap"
	"code.google.com/p/go.text/transform"
	"github.com/codegangsta/martini"

//	"github.com/gregjones/httpcache"
)

type ProgrammeRequest struct {
	// JSON date string as output by Date.toJSON()
	// e.g. 2012-04-23T18:25:43.511Z
	Days []time.Time `json:"days"`
	// Hours since midnight (24 hr time)
	Hours []int `json:"hours"`
	// Channel ID string
	Channels []string `json:"channels"`
}

type ChannelRequest struct {
	// If empty, return all channels
	ChannelName string `json:"channel_name"`
	// If true, match anything containing the name (case insensitive)
	Contains bool `json:"contains"`
}

type ChannelResponse struct {
	Error    string    `json:"error"`
	Channels []Channel `json:"channels,omitempty"`
}

type APIResponse struct {
	Error string `json:"error"`
	Data  string `json:"data,omitempty"`
}

var dataList *DataList
var ResponseLimit int = 500

func InitAPI(userAgent string) error {
	dataList = &DataList{}

	dataList.Channels = append(dataList.Channels, &Channel{Id: "7a"}, &Channel{Id: "10C", DataFor: []string{"2014-03-09"}})

	/*
		t := httpcache.NewMemoryCacheTransport()
		client := http.Client{Transport: t}
		req, err := http.NewRequest("GET", DataListFile, nil)
		if err != nil {
			return err
		}
		req.Header.Set("User-Agent", userAgent)
		log.Println("Requesting datalist: " + DataListFile)
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		decoder := xml.NewDecoder(res.Body)
		decoder.CharsetReader = CharsetReader

		decoder.Decode(&dataList)
		if err != nil {
			return err
		}
	*/
	parseDataFor()
	buildChannelMap()
	//fmt.Printf("%v\n", dataList)

	return nil
}

func ChannelHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {

	apiResponse := &ChannelResponse{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apiResponse.Error = "Error reading HTTP request"
		WriteJsonRes(w, apiResponse, http.StatusBadRequest)
		return
	}

	cr := ChannelRequest{}

	if len(body) > 0 {
		err = json.Unmarshal(body, &cr)
		if err != nil {
			apiResponse.Error = "Invalid JSON supplied, error: " + err.Error()
			WriteJsonRes(w, apiResponse, http.StatusBadRequest)
			return
		}
	}

	log.Printf("Channel Request: %+v\n", cr)

	if cr.ChannelName == "" {
		if len(dataList.ChannelMap) < ResponseLimit {
			WriteJsonRes(w, dataList.ChannelMap, http.StatusOK)
		} else {
			apiResponse.Error = fmt.Sprintf("Response count exceeded limit %d\n", ResponseLimit)
			WriteJsonRes(w, apiResponse, http.StatusBadRequest)
			return

		}
	} else if cr.Contains {
		channels := make([]*Channel, 0)
		for k, v := range dataList.ChannelMap {
			if strings.Contains(strings.ToLower(k), strings.ToLower(cr.ChannelName)) {
				channels = append(channels, v)

			}
		}
		WriteJsonRes(w, channels, http.StatusOK)
	} else {
		WriteJsonRes(w, dataList.ChannelMap[cr.ChannelName], http.StatusOK)
	}
}

func ProgrammeHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {

	apiResponse := &APIResponse{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apiResponse.Error = "Error reading HTTP request"
		WriteJsonRes(w, apiResponse, http.StatusBadRequest)
		return
	}

	pr := ProgrammeRequest{}
	err = json.Unmarshal(body, &pr)
	if err != nil {
		apiResponse.Error = "Invalid JSON supplied, error: " + err.Error()
		WriteJsonRes(w, apiResponse, http.StatusBadRequest)
		return
	}

	for _, channelr := range pr.Channels {
		log.Printf("channelr: %v\n", channelr)
		channel := dataList.ChannelMap[channelr]
		log.Printf("channel: %v\n", channel)
		log.Printf("channel map: %v\n", channel.DataForT)
		if channel != nil {
			for _, dayr := range pr.Days {
		log.Printf("dayr: %v\n", dayr)
				dayrMidnight := time.Date(dayr.Year(), dayr.Month(), dayr.Hour(), 0, 0, 0, 0, dayr.Location())
		log.Printf("dayrMidnight: %v\n", dayrMidnight)
				day := channel.DataForT[dayrMidnight]
				if day {
					filename := channel.Id + "_" + dayrMidnight.Format("2006-01-02") + ".xml.gz"
						log.Printf("Request filename: %s\n", filename)
				}
			}
		}
	}

}

func CharsetReader(charset string, input io.Reader) (io.Reader, error) {
	// Windows-1252 is a superset of ISO-8859-1.
	if charset == "iso-8859-1" {
		return transform.NewReader(input, charmap.Windows1252.NewDecoder()), nil
	}
	return nil, fmt.Errorf("unsupported charset: %q", charset)
}

func WriteJsonRes(w http.ResponseWriter, obj interface{}, statusCode int) {
	json, err := json.Marshal(&obj)
	if err != nil {
		json = []byte("{\"statusCode\": 500, \"error\": \"" + err.Error() + "\"}")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(json)
}
