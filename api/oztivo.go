package api

import (
//	"fmt"
//	"time"
)

type Channel struct {
	Id          string `xml:"id,attr"`
	DisplayName struct {
		Name string `xml:",innerxml"`
		Lang string `xml:"lang,attr"`
	} `xml:"display-name"`
	BaseURL []string `xml:"base-url"`
	dataFor []string `xml:"datafor"`
	//	DataFor_  []time.Time `xml:"-"`
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

/*
func parseDataFor() {
	fmt.Printf("DataFor enter\n")
	for _, channel := range dataList.Channels {
		fmt.Printf("DataFor channel %+v\n", channel)
		for _, df := range channel.dataFor {
			fmt.Printf("1DataFor: %v\n", df)
			t, err := time.Parse("2006-01-02", df)
			if err == nil {
				fmt.Printf("2DataFor: %v\n", t)
				channel.DataFor_ = append(channel.DataFor_, t)
			}
		}
	}
}
*/

func buildChannelMap() {
	dataList.ChannelMap = make(map[string]*Channel)
	for _, channel := range dataList.Channels {
		dataList.ChannelMap[channel.Id] = channel
	}
}
