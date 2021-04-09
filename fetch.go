package main

import (
	"bufio"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Document struct {
	Title   string
	Content string
	Url     string
}

type Fetcher interface {
	Fetch(targetUrl string) (doc *Document, urls []string, err error)
}

type DefaultFetcher struct {
}

func newDefaultFetcher() DefaultFetcher {
	return DefaultFetcher{}
}

func (fetcher DefaultFetcher) Fetch(targetUrl string) (doc *Document, urls []string, err error) {
	resp, err := http.Get(targetUrl)
	if err != nil {
		return nil, nil, err
	}
	reader := bufio.NewReader(resp.Body)
	gdoc, err := goquery.NewDocumentFromReader(reader)
	defer resp.Body.Close()
	if err != nil {
		return nil, nil, err
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
	return document, next, nil
}
