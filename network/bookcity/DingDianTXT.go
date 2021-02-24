package bookcity

import (
	"Novel/entity"
	"Novel/network"
	"Novel/utils"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

//顶点txt书城的小说资源
type DingDianTXT struct {
	Host string
}

func (DingDianTXT) RequestSearchHTML(key string) (string, error) {
	urlStr := "https://www.txtbook.net/search.php"
	urlValues := url.Values{}
	urlValues.Add("q", key)
	html, err := network.GetRequest(urlStr, urlValues)
	return html, err
}

func (receiver DingDianTXT) SourceURLHTMLParse(html, title, author string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	var sourceURL string
	//fmt.Println(doc.Find(".result-item").First().Text())
	doc.Find(".result-item.result-game-item").Each(func(i int, s *goquery.Selection) {
		searchTitle := s.Find(".result-game-item-detail h3 a span").Text()
		searchAuthor := s.Find(".result-game-item-detail .result-game-item-info p span").Next().First().Text()
		if title == searchTitle && author == searchAuthor {
			if source, exists := s.Find(".result-game-item-detail h3 a").Attr("href"); exists {
				sourceURL = receiver.Host + source
			}
			return
		}
	})
	if sourceURL == "" {
		return "", network.Errors{Code: -201, Msg: network.ErrorCode[-201]}
	}
	return sourceURL, nil
}

func (receiver DingDianTXT) DirectoryHTMLParse(html, sourceURL string) ([]entity.ChapterEntity, error) {
	if ret, err := utils.ConvertGBKToUTF8(html); err != nil {
		return nil, err
	} else {
		html = ret
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	var chapterList []entity.ChapterEntity
	doc.Find("#list dl dd").Each(func(i int, s *goquery.Selection) {
		chapter := entity.ChapterEntity{}
		chapter.Rank = uint16(i + 1)
		chapter.ChapterName = s.Find("a").Text()
		if urlStr, exists := s.Find("a").Attr("href"); exists {
			chapter.ChapterURL = receiver.Host + urlStr
		}
		chapterList = append(chapterList, chapter)
	})
	if chapterList == nil {
		return nil, network.Errors{Code: -202, Msg: network.ErrorCode[-202]}
	}
	return chapterList, nil
}

func (receiver DingDianTXT) ContentHTMLParse(html string) (string, error) {
	if ret, err := utils.ConvertGBKToUTF8(html); err != nil {
		return "", err
	} else {
		html = ret
	}
	return "", nil
}
