package main

import (
	"bufio"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type Document struct {
	Title   string
	Content string
	Url     string
}

type Fetcher interface {
	Fetch(targetUrl string) (doc *Document, urls []string, success bool, err error)
}

type DefaultFetcher struct {
}

func newDefaultFetcher() DefaultFetcher {
	return DefaultFetcher{}
}

func getCharset(header http.Header) string {
	if header != nil {
		items := strings.Split(header.Get("Content-Type"), ";")
		for _, item := range items {
			item = strings.TrimSpace(item)
			if strings.HasPrefix(item, "charset=") {
				return strings.ToLower(item[8:])
			}
		}
	}

	return "utf-8"
}

func getDecoder(charset string) (bool, transform.Transformer) {

	switch charset {
	case "gbk":
		return true, simplifiedchinese.GBK.NewDecoder()
	case "gb2312":
		fallthrough
	case "gb18030":
		return true, simplifiedchinese.GB18030.NewDecoder()
	case "hzgb2312":
		return true, simplifiedchinese.HZGB2312.NewDecoder()
	case "big-5":
	case "big5":
		return true, traditionalchinese.Big5.NewDecoder()
	}
	return false, nil
}

func (f DefaultFetcher) Fetch(targetUrl string) (doc *Document, urls []string, success bool, err error) {
	resp, err := http.Get(targetUrl)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, nil, false, nil
	}
	if err != nil {
		return nil, nil, false, err
	}
	contentCharset := getCharset(resp.Header)
	var reader io.Reader

	reader = bufio.NewReader(resp.Body)
	defer resp.Body.Close()
	if contentCharset != "utf-8" && contentCharset != "" {
		hasDecoder, decoder := getDecoder(contentCharset)
		if hasDecoder {
			reader = transform.NewReader(reader, decoder)
		}
	}
	gdoc, err := goquery.NewDocumentFromReader(reader)
	defer resp.Body.Close()
	if err != nil {
		return nil, nil, false, err
	}

	aNodes := gdoc.Find("a")
	next := make([]string, aNodes.Length())
	aNodes.Each(func(i int, s *goquery.Selection) {
		val, exist := s.Attr("href")
		if exist {
			next[i] = val
		}
	})
	document := new(Document)
	document.Url = targetUrl
	document.Title = strings.TrimSpace(gdoc.Find("title").Text())
	document.Content = gdoc.Text()
	return document, next, true, nil
}
