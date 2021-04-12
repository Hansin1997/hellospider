package core

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type Document struct {
	Title          string
	Description    string
	Content        string
	Keywords       string
	Url            string
	Host           string
	FetchAt        time.Time
	ResponseHeader http.Header
	RequestHeader  http.Header
}

type Fetcher interface {
	Fetch(targetUrl string) (doc *Document, urls []string, success bool, err error)
}

type DefaultFetcher struct {
	accepts    []string
	userAgents []string
}

func NewDefaultFetcher(accepts []string, userAgents []string) DefaultFetcher {
	return DefaultFetcher{accepts, userAgents}
}

func checkContentType(content_type string, accepts []string) (isAccept bool, conentCharset string) {
	items := strings.Split(content_type, ";")
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
	return isAccept, ""
}

func getDecoder(charset string) (bool, transform.Transformer) {
	charset = strings.ToLower(charset)
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

func shouldKeep(char rune) bool {
	switch char {
	case ' ':
		fallthrough
	case '\r':
		fallthrough
	case '\n':
		fallthrough
	case '\t':
		fallthrough
	case '\v':
		fallthrough
	case '\f':
		return true
	default:
		return false
	}
}

func trimContent(input string) (output string) {
	bder := strings.Builder{}
	lastFlag := true
	for _, c := range input {
		if shouldKeep(c) {
			if !lastFlag {
				bder.WriteRune(' ')
				lastFlag = true
			}
			continue
		}
		if lastFlag {
			lastFlag = false
		}
		bder.WriteRune(c)
	}
	return bder.String()
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
	if requ.Body != nil {
		defer requ.Body.Close()
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, nil, false, nil
	}

	isAccept, contentCharset := checkContentType(resp.Header.Get("Content-Type"), f.accepts)
	if !isAccept {
		return nil, nil, false, nil
	}

	var gdoc *goquery.Document

	if contentCharset == "" {
		// 如果 Header 中没有 charset ，则需要一次读完内容，再从 Html 元素中判断。

		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, nil, false, err
		}

		gdoc, err = goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))

		if err != nil {
			return nil, nil, false, err
		}

		selections := gdoc.Find("meta[charset]")
		if selections.Length() > 0 {
			cs, exist := selections.First().Attr("charset")
			if exist && cs != "" && !strings.EqualFold(cs, contentCharset) {
				contentCharset = strings.ToLower(cs)
			}
		} else {
			selections = gdoc.Find("meta[http-equiv='Content-Type']")
			if selections.Length() > 0 {
				cs, exist := selections.First().Attr("content")
				if exist && cs != "" {
					items := strings.Split(cs, ";")
					for _, item := range items {
						item = strings.ToLower(strings.TrimSpace(item))
						if strings.HasPrefix(item, "charset=") {
							contentCharset = item[8:]
							break
						}
					}
				}
			}
		}

		if contentCharset != "utf-8" && contentCharset != "" {
			hasDecoder, decoder := getDecoder(contentCharset)
			if hasDecoder {
				reader := transform.NewReader(bytes.NewReader(bodyBytes), decoder)
				gdoc, err = goquery.NewDocumentFromReader(reader)
				if err != nil {
					return nil, nil, false, err
				}
			}
		}
	} else {
		// Header 中有 charset ，直接读 Body 流。
		var reader io.Reader = bufio.NewReader(resp.Body)

		// 如果字符编码不是 utf-8 ，进行转换
		if contentCharset != "utf-8" {
			hasDecoder, decoder := getDecoder(contentCharset)
			if hasDecoder {
				reader = transform.NewReader(reader, decoder)
			}
		}

		gdoc, err = goquery.NewDocumentFromReader(reader)

		if err != nil {
			return nil, nil, false, err
		}
	}

	aNodes := gdoc.Find("a")
	next := make([]string, aNodes.Length())
	aNodes.Each(func(i int, s *goquery.Selection) {
		val, exist := s.Attr("href")
		if exist {
			next[i] = val
		}
	})
	gdoc.Find("script").Remove()
	gdoc.Find("style").Remove()
	document := new(Document)
	document.Url = targetUrl
	document.Title = strings.TrimSpace(gdoc.Find("title").Text())
	document.Content = trimContent(gdoc.Find("body").Text())
	document.Description = strings.TrimSpace(gdoc.Find("meta[name='description']").AttrOr("content", ""))
	document.Keywords = strings.TrimSpace(gdoc.Find("meta[name='description']").AttrOr("keywords", ""))

	document.RequestHeader = resp.Request.Header
	document.ResponseHeader = resp.Header
	document.FetchAt = time.Now()
	document.Host = resp.Request.URL.Host
	return document, next, true, nil
}
