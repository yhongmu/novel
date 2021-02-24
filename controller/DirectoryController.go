package controller

import (
	"Novel/dao"
	"Novel/entity"
	"Novel/network"
	"Novel/network/bookcity"
	"net/http"
)

var Directory *directoryController

func GetDirectoryInstance() *directoryController {
	if Directory == nil {
		Directory = &directoryController{}
	}
	return Directory
}

type directoryController struct {
}

func (p *directoryController) Router(router *network.RouterHandler) {
	router.Router("/novel/directory", p.getDirectory)
}

func (p *directoryController) getDirectory(w http.ResponseWriter, r *http.Request)  {
	novelID := r.PostFormValue("id")
	if novelID == "" {
		network.ResultFail(w, -101, network.ErrorCode[-101])
		return
	}
	directoryDao := dao.DirectoryDao{TableName: "dir_" + novelID}
	if err := directoryDao.CreateNovelDirectoryTable(); err != nil {
		network.ErrorShow(w, err)
		return
	}
	source, sourceURL, err := dao.GetInfoDaoInstance().QuerySourceByNovelID(novelID)
	if err != nil {
		network.ErrorShow(w, err)
		return
	}
	chapterList, err := p.getDirectoryRequest(source, sourceURL)
	if err != nil {
		network.ErrorShow(w, err)
		return
	}
	if err := directoryDao.InsertAllData(chapterList); err != nil {
		network.ErrorShow(w, err)
		return
	}
	network.ResultOK(w, 0, "小说目录获取成功！", chapterList)
}

func (p *directoryController) getDirectoryRequest(source, sourceURL string) ([]entity.ChapterEntity, error) {
	handler := bookcity.SelectResources(source)
	html, err := network.GetRequest(sourceURL, nil)
	if err != nil {
		return nil, err
	}
	if chapterList, err := handler.DirectoryHTMLParse(html, sourceURL); err != nil {
		return nil, err
	} else {
		return chapterList, nil
	}

}
