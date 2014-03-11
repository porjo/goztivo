package api

import (
	"encoding/json"
	"encoding/xml"
	"errors"
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
	"github.com/gregjones/httpcache"
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

type HttpError struct {
	msg        string
	statusCode int
}

func (h HttpError) Error() string {
	return h.msg
}

type APIResponse struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type FileRequest struct {
	Channel  string
	Filename string
}

var (
	dataList      *DataList
	ResponseLimit int = 500
	transport     *httpcache.Transport
	userAgent     string
)

func InitAPI(userAgentIn string) error {

	userAgent = userAgentIn
	dataList = &DataList{}

	//		dataList.Channels = append(dataList.Channels, &Channel{Id: "7A", DataFor: []string{"2014-03-08", "2014-03-09", "2014-03-10"}}, &Channel{Id: "10C", DataFor: []string{"2014-03-08", "2014-03-09", "2014-03-10"}}, &Channel{Id: "9B", DataFor: []string{"2014-03-08", "2014-03-09", "2014-03-10"}})

	transport = httpcache.NewMemoryCacheTransport()
	client := transport.Client()
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

	dataList.parseDataFor()
	dataList.buildChannelMap()
	//fmt.Printf("%v\n", dataList.ChannelMap)

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

	fileRequests, err := pr.buildFileList()
	if err != nil {
		apiResponse.Error = err.Error()
		WriteJsonRes(w, apiResponse, http.StatusBadRequest)
		return
	}

	if len(fileRequests) == 0 {
		apiResponse.Error = "No Results found"
		WriteJsonRes(w, apiResponse, http.StatusOK)
		return
	}

	channelDays, err := fetchChannelDays(fileRequests)
	if err != nil {
		apiResponse.Error = err.Error()
		if e, ok := err.(HttpError); ok {
			WriteJsonRes(w, apiResponse, e.statusCode)
		} else {
			WriteJsonRes(w, apiResponse, http.StatusInternalServerError)
		}
		return
	}

	for _, channelDay := range channelDays {
		channelDay.parseStopStart()
	}

	apiResponse.Data = &channelDays
	WriteJsonRes(w, apiResponse, http.StatusOK)
	return
}

// Fetch filenames
func fetchChannelDays(fileRequests []FileRequest) (channelDays []ChannelDay, err error) {
	channelDays = make([]ChannelDay, 0)
	client := transport.Client()
	for _, fileRequest := range fileRequests {

		var req *http.Request
		var res *http.Response

		url := BaseURL + fileRequest.Filename
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
		req.Header.Set("User-Agent", userAgent)
		log.Println("Fetching URL: " + url)
		res, err = client.Do(req)
		if err != nil {
			return
		}
		if res.StatusCode != 200 {
			errMsg := fmt.Sprintf("Remote server returned status code: %d when fetching '%s'", res.StatusCode, url)
			err = HttpError{errMsg, 502}
			return
		}
		decoder := xml.NewDecoder(res.Body)
		decoder.CharsetReader = CharsetReader
		channelDay := ChannelDay{}
		channelDay.Channel = fileRequest.Channel
		err = decoder.Decode(&channelDay)
		if err != nil {
			res.Body.Close()
			return
		}
		channelDays = append(channelDays, channelDay)
		res.Body.Close()
	}
	return
}

func (pr ProgrammeRequest) buildFileList() (fileRequests []FileRequest, err error) {
	fileRequests = make([]FileRequest, 0)
	for _, channelr := range pr.Channels {
		channel := dataList.ChannelMap[channelr]
		if channel != nil {
			for _, dayr := range pr.Days {
				dayrMidnight := time.Date(dayr.Year(), dayr.Month(), dayr.Day(), 0, 0, 0, 0, dayr.Location())
				day := channel.DataForT[dayrMidnight]
				if day {
					filename := channel.Id + "_" + dayrMidnight.Format("2006-01-02") + ".xml.gz"
					fileRequest := FileRequest{channelr, filename}
					fileRequests = append(fileRequests, fileRequest)
				}
			}
		} else {
			err = errors.New("Channel " + channelr + " not found")
		}
	}
	return
}

func CharsetReader(charset string, input io.Reader) (io.Reader, error) {
	// Windows-1252 is a superset of ISO-8859-1.
	if strings.ToLower(charset) == "iso-8859-1" {
		return transform.NewReader(input, charmap.Windows1252.NewDecoder()), nil
	}
	return nil, fmt.Errorf("unsupported charset: %q", charset)
}

func WriteJsonRes(w http.ResponseWriter, obj interface{}, statusCode int) {
	json, err := json.Marshal(&obj)
	if err != nil {
		json = []byte("{\"statusCode\": 500, \"error\": \"" + err.Error() + "\"}")
	}

	// JSON Vulnerability Protection for AngularJS
	ngPrefix := []byte(")]}',\n")
	ngPrefix = append(ngPrefix, json...)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(ngPrefix)
}
