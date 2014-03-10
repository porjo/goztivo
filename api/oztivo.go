package api

import (
	"log"
	"time"
)

type Channel struct {
	Id          string `xml:"id,attr"`
	DisplayName struct {
		Text string `xml:",innerxml"`
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
	StartTime   string    `xml:"start,attr" json:"-"`
	StopTime    string    `xml:"stop,attr" json:"-"`
	StartTimeJ  time.Time `xml:"-" json:"start_time"`
	StopTimeJ   time.Time `xml:"-" json:"stop_time"`
	Title       string    `xml:"title"`
	SubTitle    string    `xml:"sub-title" json:"subtitle"`
	Description string    `xml:"desc" json:"description"`
	Credits     []struct {
		Actor string `xml:"actor"`
	} `xml:"credits"`
	Category []string `xml:"category"`
	Rating   []struct {
		Value string `xml:"value"`
	} `xml:"rating"`
	StarRating []struct {
		Value string `xml:"value"`
	} `xml:"star-rating" json:"star_rating"`
}

type ChannelDay struct {
	Programmes []*Programme `xml:"programme"`
}

const (
	BaseURL      = "http://xml.oztivo.net/xmltv/"
	DataListFile = BaseURL + "datalist.xml.gz"
)

// Convert time strings to time.Time
func (d *DataList) parseDataFor() {
	for _, channel := range d.Channels {
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
		log.Printf("Add to map %s\n", channel)
		d.ChannelMap[channel.Id] = channel
	}
}
