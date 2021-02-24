package controller

import (
	"Novel/dao"
	"Novel/entity"
	"Novel/network"
	"Novel/network/bookcity"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

var Content *contentController
var wg sync.WaitGroup

func GetContentInstance() *contentController {
	if Content == nil {
		Content = &contentController{}
	}
	return Content
}

type contentController struct {
}

func (p *contentController) Router(router *network.RouterHandler) {
	router.Router("/novel/content", p.getContent)
}

func (p contentController) getContent(w http.ResponseWriter, r *http.Request) {
	novelID := r.PostFormValue("id")
	rank := r.PostFormValue("rank")
	if novelID == "" || rank == "" {
		network.ResultFail(w, -101, network.ErrorCode[-101])
		return
	}
	source := novelID[0: 3]
	rankUint, err := strconv.ParseUint(rank, 10, 16)
	if err != nil {
		network.ErrorShow(w, err)
		return
	}
	contentDao := dao.ContentDao{TableName: "text_" + novelID}
	//创建小说内容表
	if err := contentDao.CreateNovelContentTable(); err != nil {
		network.ErrorShow(w, err)
		return
	}
	var chapterList []entity.ChapterEntity
	//判断该rank的小说章节是否存储在数据库中
	if contentDao.HasContent(rankUint) {
		directoryDao := dao.DirectoryDao{TableName: "dir_" + novelID}
		//去小说目录表获取URL
		if chapterList, err = directoryDao.SelectDIRByRank(rankUint); err != nil {
			//没获取到url到情况
			network.ErrorShow(w, err)
			return
		} else {
			//获取到url后，去爬取内容数据
			if chapterList, err = p.getContentRequest(source, chapterList); err != nil {
				network.ErrorShow(w, err)
				return
			} else {
				if err = contentDao.InsertAllContent(chapterList); err != nil {
					network.ErrorShow(w, err)
					return
				}
				network.ResultOK(w, 0, "小说正文获取成功！", chapterList)
			}
		}
	} else {
		//获取数据库的小说正文内容，并返回
		if chapterList, err = contentDao.SelectContentByRank(rankUint); err != nil {
			network.ErrorShow(w, err)
			return
		} else {
			network.ResultOK(w, 0, "小说正文获取成功！", chapterList)
		}
	}
}

func (p contentController) getContentRequest(source string, chapterList []entity.ChapterEntity) ([]entity.ChapterEntity, error) {
	handler := bookcity.SelectResources(source)
	var list []entity.ChapterEntity
	taskIndex := make(chan int, entity.CHAPTER_RANGE)
	taskHTML := make(chan string, entity.CHAPTER_RANGE)
	wg.Add(int(entity.CHAPTER_RANGE))
	for index, chapter := range chapterList {
		go p.goWorkRequest(chapter.ChapterURL, index, taskIndex, taskHTML)
	}
	for range chapterList {
		index, ok := <-taskIndex
		if !ok {
			return nil, network.Errors{Code: -401, Msg: network.ErrorCode[-401]}
		}
		html := <-taskHTML
		content, err := handler.ContentHTMLParse(html)
		if err != nil {
			return nil, err
		}
		entity.ContentPaging(content, chapterList[index], &list)
	}
	close(taskIndex)
	close(taskHTML)
	// 等待小说内容爬取
	wg.Wait()
	sort.Slice(list, func(i, j int) bool {
		return list[i].Rank < list[j].Rank
	})
	return list, nil
}

func (p contentController) goWorkRequest(url string, index int, taskIndex chan int, taskHTML chan string) {
	defer wg.Done()
	html, err := network.GetRequest(url, nil)
	if err != nil {
		return
	}
	var mutex sync.Mutex
	mutex.Lock()
	{
		taskIndex <- index
		taskHTML <- html
	}
	mutex.Unlock()
}
