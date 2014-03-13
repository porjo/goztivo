package api

import (
	"log"
	"sync"
	"time"
)

type TimeList []time.Time

type Channel struct {
	Id          string `xml:"id,attr" json:"id"`
	DisplayName struct {
		Text string `xml:",innerxml" json:"text"`
		Lang string `xml:"lang,attr" json:"lang"`
	} `xml:"display-name" json:"display_name"`
	BaseURL  []string `xml:"base-url" json:"-"`
	DataFor  []string `xml:"datafor" json:"-"`
	DataForT TimeList `xml:"-" json:"-"`
}

type DataList struct {
	Channels []*Channel `xml:"channel"`
	// Maps channel id to channel
	ChannelMap map[string]*Channel `xml:"-"`
	Mutex       *sync.Mutex
}

type Programme struct {
	StartTime   string    `xml:"start,attr" json:"-"`
	StopTime    string    `xml:"stop,attr" json:"-"`
	StartTimeJ  time.Time `xml:"-" json:"start_time"`
	StopTimeJ   time.Time `xml:"-" json:"stop_time"`
	Title       string    `xml:"title" json:"title"`
	SubTitle    string    `xml:"sub-title" json:"subtitle"`
	Description string    `xml:"desc" json:"description"`
	Credits     []struct {
		Actor string `xml:"actor" json:"actor,omitempty"`
	} `xml:"credits" json:"credits,omitempty"`
	Category []string `xml:"category" json:"category,omitempty"`
	Rating   []struct {
		Value string `xml:"value" json:"value,omitempty"`
	} `xml:"rating" json:"rating,omitempty"`
	StarRating []struct {
		Value string `xml:"value" json"value,omitempty"`
	} `xml:"star-rating" json:"star_rating,omitempty"`
}

type ChannelDay struct {
	ChannelId  string       `xml:"-" json:"id"`
	Programmes []*Programme `xml:"programme" json:"programme"`
}

const (
	BaseURL      = "http://xml.oztivo.net/xmltv/"
	DataListFile = BaseURL + "datalist.xml.gz"
)

// Convert time strings to time.Time
func (d *DataList) parseDataFor() {
	for _, channel := range d.Channels {
		channel.DataForT = nil
		for _, df := range channel.DataFor {
			// Dates are relative to AEST (I presume)
			tz, err := time.Parse("2006-01-02", df)
			loc, err := time.LoadLocation("Australia/Sydney")
			if err != nil {
				panic("Could not load timezone location")
			}
			t := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, loc)
			if err != nil {
				log.Println(df, err)
			} else {
				channel.DataForT = append(channel.DataForT, t)
			}
		}

		// Clear slice
		channel.DataFor = nil
	}
}

// Convert time strings to time.Time
func (c ChannelDay) parseStopStart() {
	for _, programme := range c.Programmes {
		var err error
		programme.StartTimeJ, err = time.Parse("20060102150405 -0700", programme.StartTime)
		if err != nil {
			log.Println(err)
		}
		programme.StopTimeJ, err = time.Parse("20060102150405 -0700", programme.StopTime)
		if err != nil {
			log.Println(err)
		}
	}
}

// Convert slice of channels to map of channels
func (d *DataList) buildChannelMap() {
	d.ChannelMap = make(map[string]*Channel)
	for _, channel := range d.Channels {
		//log.Printf("Add to map %s\n", channel)
		d.ChannelMap[channel.Id] = channel
	}
}
