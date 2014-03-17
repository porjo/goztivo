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
	"sync"
	"time"

	"code.google.com/p/go.text/encoding/charmap"
	"code.google.com/p/go.text/transform"
	"github.com/codegangsta/martini"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

type ProgrammeRequest struct {
	// JSON date string as output by Date.toJSON()
	// e.g. 2012-04-23T18:25:43.511Z
	Days []time.Time `json:"days"`

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
	Date     time.Time
}

type OzClient struct {
	Client    *http.Client
	UserAgent string
	// Lock client access to enforce one request
	// at a time to the upstream server
	Mutex *sync.Mutex
	// Keep track of when we last fetched, to keep
	// our upstream requests spaced apart
	LastFetch time.Time
}

var (
	dataList      *DataList
	ResponseLimit int = 500
	transport     *httpcache.Transport
	ozclient      OzClient
)

func InitAPI(userAgentIn string) {

	ozclient.UserAgent = userAgentIn

	tempDir, err := ioutil.TempDir("", "goztivo_")
	if err != nil {
		panic(err)
	}
	cache := diskcache.New(tempDir)
	transport = httpcache.NewTransport(cache)

	ozclient = OzClient{}
	ozclient.Client = transport.Client()
	ozclient.Mutex = &sync.Mutex{}
	ozclient.LastFetch = time.Now()

	dataList = &DataList{}
	dataList.Mutex = &sync.Mutex{}

	go func() {

		for {
			dataList.Mutex.Lock()
			err := GetDataList()
			dataList.Mutex.Unlock()
			if err != nil {
				log.Println(err)
				return
			}

			time.Sleep(time.Hour * 24)
		}

	}()
}

func GetDataList() error {

	dataList.Channels = nil
	dataList.ChannelMap = nil

	req, err := http.NewRequest("GET", DataListFile, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", ozclient.UserAgent)
	log.Println("Requesting datalist: " + DataListFile)
	ozclient.Mutex.Lock()
	res, err := ozclient.Client.Do(req)
	ozclient.Mutex.Unlock()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := xml.NewDecoder(res.Body)
	decoder.CharsetReader = charsetReader

	decoder.Decode(&dataList)
	if err != nil {
		return err
	}

	dataList.parseDataFor()
	dataList.buildChannelMap()

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

	log.Printf("Channel request: %+v\n", cr)

	dataList.Mutex.Lock()
	defer dataList.Mutex.Unlock()

	if cr.ChannelName == "" {
		if len(dataList.ChannelMap) < ResponseLimit {
			WriteJsonRes(w, dataList.ChannelMap, http.StatusOK)
		} else {
			apiResponse.Error = fmt.Sprintf("Response count exceeded limit %d\n", ResponseLimit)
			WriteJsonRes(w, apiResponse, http.StatusBadRequest)
			return

		}
	} else if cr.Contains {
		var channels []*Channel
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

	log.Printf("Programme request: %+v\n", pr)

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

	apiResponse.Data = &channelDays
	WriteJsonRes(w, apiResponse, http.StatusOK)
	return
}

// Fetch remote files
func fetchChannelDays(fileRequests []FileRequest) (channelDays []ChannelDay, err error) {
	for _, fileRequest := range fileRequests {

		var req *http.Request
		var res *http.Response

		url := BaseURL + fileRequest.Filename
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
		req.Header.Set("User-Agent", ozclient.UserAgent)

		// First fetch is to simply use whatever is in cache. This
		// lets us respond to the client quicker by not having to wait
		// for remote server to confirm freshness of file.
		// We then do a second fetch in the background to update
		// the cache for subsequent requests
		req.Header.Set("Cache-Control", "max-stale=99999")
		log.Println("Fetching URL: " + url)
		res, err = ozclient.Client.Do(req)
		if err != nil {
			return
		}
		if res.StatusCode != 200 {
			errMsg := fmt.Sprintf("Remote server returned status code: %d when fetching '%s'", res.StatusCode, url)
			err = HttpError{errMsg, 502}
			return
		}

		go func() {
			ozclient.Mutex.Lock()
			if time.Since(ozclient.LastFetch) < time.Second {
				log.Println("Fetching too quickly, sleeping...")
				//Sleep between successive file gets (as per Oztivo usage policy)
				time.Sleep(time.Second)
			}
			req.Header.Set("Cache-Control", "max-age=0")
			res, err = ozclient.Client.Do(req)
			ozclient.Mutex.Unlock()
			ozclient.LastFetch = time.Now()

			if err != nil {
				return
			}
			if res.StatusCode != 200 {
				errMsg := fmt.Sprintf("Remote server returned status code: %d when fetching '%s'", res.StatusCode, url)
				err = HttpError{errMsg, 502}
				return
			}
		}()

		decoder := xml.NewDecoder(res.Body)
		decoder.CharsetReader = charsetReader
		channelDay := ChannelDay{}
		channelDay.ChannelId = fileRequest.Channel
		err = decoder.Decode(&channelDay)
		if err != nil {
			res.Body.Close()
			return
		}
		channelDay.parseStopStart()
		channelDay.Date = fileRequest.Date
		channelDays = append(channelDays, channelDay)
		res.Body.Close()
	}
	return
}

func (pr ProgrammeRequest) buildFileList() (fileRequests []FileRequest, err error) {
	for _, channelr := range pr.Channels {
		dataList.Mutex.Lock()
		channel := dataList.ChannelMap[channelr]
		dataList.Mutex.Unlock()
		if channel != nil {
			for _, tz := range pr.Days {
				if channel.DataForT.contains(tz) {
					loc, err := time.LoadLocation("Australia/Sydney")
					if err != nil {
						panic("Could not load timezone location")
					}
					t := time.Date(tz.Year(), tz.Month(), tz.Day(), tz.Hour(), tz.Minute(), tz.Second(), 0, loc)
					filename := channel.Id + "_" + t.Format("2006-01-02") + ".xml.gz"
					fileRequest := FileRequest{channelr, filename, tz}
					fileRequests = append(fileRequests, fileRequest)
				}
			}
		} else {
			err = errors.New("Channel " + channelr + " not found")
		}

	}
	return
}

func (tlist TimeList) contains(t time.Time) bool {
	for _, x := range tlist {
		startOfDay := time.Date(x.Year(), x.Month(), x.Day(), 0, 0, 0, 0, x.Location())
		endOfDay := time.Date(x.Year(), x.Month(), x.Day(), 23, 59, 59, 99999, x.Location())

		if (t.After(startOfDay) || t.Equal(startOfDay)) && (t.Before(endOfDay) || t.Equal(endOfDay)) {
			return true
		}
	}
	return false
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
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
