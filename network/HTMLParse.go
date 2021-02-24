package network

import (
	"Novel/entity"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"regexp"
)



func SearchResultHTMLParse(r io.Reader) (interface{}, error) {
	var novelInfoList []entity.NovelInfoEntity
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// 解析书籍图片url的函数
	imgParse := func(s *goquery.Selection, selector string) string {
		imgURL, _ :=s.Find(selector).First().Attr("src")
		re := regexp.MustCompile(`\w.*`)
		return re.FindString(imgURL)
	}
	// 解析书籍其他信息的函数
	detailsParse := func(s *goquery.Selection, selector string) (title, author, category, status, des, update string){
		s = s.Find(selector).First()
		title = s.Find("h4 a").Text()
		author = s.Find(".author a[data-eid=qd_S06]").Text()
		category = s.Find(".author a[data-eid=qd_S07]").Text()
		status = s.Find(".author span").Text()
		des = s.Find("p.intro").Text()
		// 保证小说描述字段的前后没有空格和换行符
		re := regexp.MustCompile(`(?s)([^\s\n].*?)[\s\n]*$`)
		match := re.FindAllStringSubmatch(des, -1)
		des = match[0][1]
		update = s.Find(".update a").Text()
		return title, author, category, status, des, update
	}
	doc.Find(".res-book-item").Each(func(i int, s *goquery.Selection) {
		info := entity.NovelInfoEntity{}
		info.ID, _ = s.Attr("data-bid")
		info.BookImg = imgParse(s, ".book-img-box a img")
		info.Title, info.Author, info.Category, info.Status, info.Description, info.Update = detailsParse(s, ".book-mid-info")
		novelInfoList = append(novelInfoList, info)
	})
	return novelInfoList, nil
}
