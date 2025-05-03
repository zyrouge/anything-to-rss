package rss

import (
	"encoding/xml"
	"net/http"
	"time"
)

type RssXml struct {
	XMLName xml.Name      `xml:"rss"`
	Version int           `xml:"version,attr"`
	Channel RssXmlChannel `xml:"channel"`
}

type RssXmlChannel struct {
	Title       string              `xml:"title"`
	Link        string              `xml:"link"`
	Description string              `xml:"description"`
	Items       []RssXmlChannelItem `xml:"item"`
}

type RssXmlChannelItem struct {
	Title       string                   `xml:"title"`
	Description string                   `xml:"description"`
	Author      string                   `xml:"author"`
	Link        string                   `xml:"link"`
	PubDate     RssXmlChannelItemPubDate `xml:"pubDate"`
}

type RssXmlChannelItemPubDate string

func (rss *RssXml) Xml() ([]byte, error) {
	bytes, err := xml.MarshalIndent(rss, "", "    ")
	if err != nil {
		return nil, err
	}
	bytesWithHeader := append([]byte(xml.Header), bytes...)
	return bytesWithHeader, nil
}

func (rss *RssXml) WriteToHttpResponseWriter(w http.ResponseWriter) {
	xml, err := rss.Xml()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	respHeaders := w.Header()
	respHeaders.Add("Content-Type", "application/rss+xml")
	w.WriteHeader(http.StatusOK)
	w.Write(xml)
}

func MakeRssXmlChannelItemPubDate(value time.Time) RssXmlChannelItemPubDate {
	return RssXmlChannelItemPubDate(value.Format(time.RFC1123))
}
