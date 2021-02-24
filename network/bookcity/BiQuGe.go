package bookcity

import (
	"Novel/network"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

//新笔趣阁书城的小说资源
type BiQuGe struct {
}

func (BiQuGe) RequestSearchHTML(key string) (string, error) {
	//新笔趣阁的搜索url
	urlStr := "http://www.xbiquge.la/modules/article/waps.php"
	urlValues := url.Values{}
	urlValues.Add("searchkey", key)
	html, err := network.PostRequest(urlStr, urlValues)
	return html, err
}

func (BiQuGe) SourceURLHTMLParse(html, title, author string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	var sourceURL string
	doc.Find("table.grid tbody tr").Each(func(i int, s *goquery.Selection) {
		//判断搜索的书名和作者是否准确
		if title == s.Find("td.even").First().Text() && author == s.Find("td.odd").Next().Text() {
			sourceURL, _ = s.Find("td.even a").First().Attr("href")
			return
		}
	})
	if sourceURL == "" {
		return "", network.Errors{Code: -201, Msg: network.ErrorCode[-201]}
	}
	return sourceURL, nil
}

