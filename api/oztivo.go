package api

import (
	"time"
	"log"
)

type Channel struct {
	Id          string `xml:"id,attr"`
	DisplayName struct {
		Name string `xml:",innerxml"`
		Lang string `xml:"lang,attr"`
	} `xml:"display-name"`
	BaseURL  []string           `xml:"base-url" json:"-"`
	DataFor  []string           `xml:"datafor" json:"-"`
	DataForT map[time.Time]bool `xml:"-" json:"-"`
}

type DataList struct {
	Channels []*Channel `xml:"channel"`
	// Maps channel id to channel
	ChannelMap map[string]*Channel `-`
}

type Programme struct {
	StartTime string `xml:"start,attr"`
	StopTime  string `xml:"stop,attr"`
	Title     string `xml:"title"`
	SubTitle  string `xml:"sub-title"`
	Desc      string `xml:"title"`
	Credits   []struct {
		Actor string `xml:"actor"`
	} `xml:"credits"`
	Category []string `xml:"category"`
	Rating   []struct {
		Value string `xml:"value"`
	} `xml:"rating"`
	StarRating []struct {
		Value string `xml:"value"`
	} `xml:"star-rating"`
}

type ChannelDay struct {
	Programmes []Programme `xml:"programme"`
}

const (
	BaseURL      = "http://xml.oztivo.net/xmltv/"
	DataListFile = BaseURL + "datalist.xml.gz"
)

// Convert time strings to time.Time
func parseDataFor() {
	for _, channel := range dataList.Channels {
		if channel.DataForT == nil {
			channel.DataForT = make(map[time.Time]bool)
		}
		for _, df := range channel.DataFor {
			t, err := time.Parse("2006-01-02", df)
			if err != nil {
				log.Println(df, err)
			} else {
				channel.DataForT[t] = true
			}
				log.Printf("df %v t %v dft %v\n", df, t, channel.DataForT)
		}

		// Clear slice
		channel.DataFor = nil
	}
}

// Convert slice of channels to map of channels
func buildChannelMap() {
	dataList.ChannelMap = make(map[string]*Channel)
	for _, channel := range dataList.Channels {
		dataList.ChannelMap[channel.Id] = channel
	}
}
