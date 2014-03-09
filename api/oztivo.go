package api

type Channel struct {
	Id          string `xml:"id,attr"`
	DisplayName struct {
		Name string `xml:",innerxml"`
		Lang string `xml:"lang,attr"`
	} `xml:"display-name"`
	BaseURL []string `xml:"base-url"`
	DataFor []string `xml:"datafor"`
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
	BaseURL = "http://xml.oztivo.net/xmltv/"
	DataListFile = BaseURL + "datalist.xml.gz"
)
