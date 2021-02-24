package bookcity

import (
	"Novel/entity"
)

// @title
//
func SelectResources(source string) ResourcesHandler {
	switch source {
	case "001":
		return ShuKu56{Host: "http://www.ppskw.com"}
	case "002":
		return DingDianTXT{Host: "https://www.txtbook.net"}
	default:
		return nil
	}
}

type ResourcesHandler interface {
	//请求对应书城的搜索接口，返回搜索的结果html
	RequestSearchHTML(string) (string, error)
	//解析搜索页面的html，获取到对应的小说目录资源url
	SourceURLHTMLParse(html, title, author string) (string, error)
	//解析小说目录的html，获取到每个章节的名字和url
	DirectoryHTMLParse(html, sourceURL string) ([]entity.ChapterEntity, error)
	//解析对应章节正文内容的html，返回小说正文内容
	ContentHTMLParse(html string) (string, error)
}
