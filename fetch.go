package main

import (
	"bufio"
	"io"
	"math/rand"
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
	accepts    []string
	userAgents []string
}

func newDefaultFetcher(accepts []string, userAgents []string) DefaultFetcher {
	return DefaultFetcher{accepts, userAgents}
}

func checkContentType(header http.Header, accepts []string) (isAccept bool, conentCharset string) {
	if header != nil {
		items := strings.Split(header.Get("Content-Type"), ";")
		if len(items) == 0 {
			return false, ""
		}
		contentType := items[0]
		for _, accept := range accepts {
			if accept == contentType {
				isAccept = true
				break
			}
		}
		if !isAccept {
			return false, ""
		}
		if len(items) > 1 {
			for _, item := range items[1:] {
				item = strings.TrimSpace(item)
				if strings.HasPrefix(item, "charset=") {
					return isAccept, strings.ToLower(item[8:])
				}
			}
		}
		return isAccept, "utf-8"
	}
	return false, ""
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

func (f DefaultFetcher) randomAgent() string {
	len := len(f.userAgents)
	if len > 0 {
		return f.userAgents[rand.Intn(len)]
	}
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36 Edg/89.0.774.75"
}

func (f DefaultFetcher) Fetch(targetUrl string) (doc *Document, urls []string, success bool, err error) {
	requ, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, nil, false, err
	}
	requ.Header.Add("User-Agent", f.randomAgent())
	resp, err := http.DefaultClient.Do(requ)
	if err != nil {
		return nil, nil, false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, nil, false, nil
	}
	isAccept, contentCharset := checkContentType(resp.Header, f.accepts)
	if !isAccept {
		return nil, nil, false, nil
	}
	var reader io.Reader

	reader = bufio.NewReader(resp.Body)
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
