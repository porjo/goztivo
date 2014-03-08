// goztivo provides a grabber for Oztivo data
// 
// Attempts to implement mostly RFC-compliant cache for http responses.
//
package main

import (
	//	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	//"time"

	"code.google.com/p/go.text/encoding/charmap"
	"code.google.com/p/go.text/transform"
	"github.com/gregjones/httpcache"
)

func main() {

	type Channel struct {
		Id          string   `xml:"id,attr"`
		DisplayName struct {
			Name string   `xml:",innerxml"`
			Lang string   `xml:"lang,attr"`
		} `xml:"display-name"`
		BaseURL     []string `xml:"base-url"`
		DataFor     []string `xml:"datafor"`
	}

	type DataList struct {
		Channels []Channel `xml:"channel"`
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

	t := httpcache.NewMemoryCacheTransport()
	client := http.Client{Transport: t}
	req, err := http.NewRequest("GET", "http://xml.oztivo.net/xmltv/datalist.xml.gz", nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	decoder := xml.NewDecoder(res.Body)
	decoder.CharsetReader = CharsetReader

	//	body, err := ioutil.ReadAll(res.Body)

	tv := DataList{}

	//err = xml.Unmarshal(body, &tv)
	decoder.Decode(&tv)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Printf("%+v\n", tv)

}

func CharsetReader(charset string, input io.Reader) (io.Reader, error) {
	// Windows-1252 is a superset of ISO-8859-1.
	if charset == "iso-8859-1" {
		return transform.NewReader(input, charmap.Windows1252.NewDecoder()), nil
	}
	return nil, fmt.Errorf("unsupported charset: %q", charset)
}
