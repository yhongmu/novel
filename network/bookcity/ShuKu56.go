package bookcity

import (
	"Novel/entity"
	"Novel/network"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

type ShuKu56 struct {
	Host string
}

func (ShuKu56) RequestSearchHTML(key string) (string, error) {
	urlStr := "http://www.ppskw.com/modules/article/search.php"
	urlValues := url.Values{}
	urlValues.Add("searchkey", key)
	html, err := network.PostRequest(urlStr, urlValues)
	return html, err
}

func (receiver ShuKu56) SourceURLHTMLParse(html, title, author string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	var sourceURL string
	doc.Find("#articlelist ul li").Each(func(i int, s *goquery.Selection) {
		searchTitle := s.Find("span.l2 a").Text()
		searchAuthor := s.Find("span.l8").Next().First().Text()
		if title == searchTitle && author == searchAuthor {
			if source, exists := s.Find("span.l2 a").Attr("href"); exists {
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

func (receiver ShuKu56) DirectoryHTMLParse(html, sourceURL string) ([]entity.ChapterEntity, error) {
	//if ret, err := utils.ConvertGBKToUTF8(html); err != nil {
	//	return nil, err
	//} else {
	//	html = ret
	//}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	var chapterList []entity.ChapterEntity
	k := 1
	doc.Find("#defaulthtml4 table tbody tr").Each(func(i int, selection *goquery.Selection) {
		selection.Find("td").Each(func(j int, s *goquery.Selection) {
			if s.Find(".dccss a").Text() == "" {
				return
			}
			chapter := entity.ChapterEntity{}
			chapter.Rank = uint16(k)
			k++
			chapter.ChapterName = s.Find(".dccss a").Text()
			if urlStr, exists := s.Find(".dccss a").Attr("href"); exists {
				chapter.ChapterURL = sourceURL + urlStr
			}
			chapterList = append(chapterList, chapter)
		})

	})
	if chapterList == nil {
		return nil, network.Errors{Code: -202, Msg: network.ErrorCode[-202]}
	}
	return chapterList, nil
}

func (receiver ShuKu56) ContentHTMLParse(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	content, err := doc.Find("#content p").Html()
	if err != nil {
		return "", err
	}
	//content = strings.Replace(content, "&nbsp", " ", -1)
	content = strings.Replace(content, "<br/>", "\n", -1)
	return content, nil
}
