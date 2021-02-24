package controller

import (
	"Novel/dao"
	"Novel/entity"
	"Novel/network"
	"Novel/network/bookcity"
	"net/http"
)

var resources *resourcesController

func GetResourcesInstance() *resourcesController {
	if resources == nil {
		resources = &resourcesController{}
	}
	return resources
}

type resourcesController struct {
}

func (p *resourcesController) Router(router *network.RouterHandler)  {
	router.Router("/novel/resources", p.getSources)
}

func (p *resourcesController) getSources(w http.ResponseWriter, r *http.Request)  {
	novelInfo := entity.NovelInfoEntity{
		ID: r.PostFormValue("id"),
		Title: r.PostFormValue("title"),
		Author: r.PostFormValue("author"),
		Category: r.PostFormValue("category"),
		Status: r.PostFormValue("status"),
		Description: r.PostFormValue("description"),
		Update: r.PostFormValue("update"),
		BookImg: r.PostFormValue("book_img"),
		Source: r.PostFormValue("source"),
	}
	if novelInfo.ID == "" || novelInfo.Title == "" || novelInfo.Author == "" {
		network.ResultFail(w, -102, network.ErrorCode[-102])
		return
	}
	//不传入source时，保证默认书源
	if novelInfo.Source == "" {
		novelInfo.Source = "001"
	}
	novelInfo.ID = novelInfo.Source + "_" + novelInfo.ID
	info, err := p.getSourcesRequest(novelInfo)
	if err != nil {
		network.ErrorShow(w, err)
		return
	}
	err = p.SourcesDBOperation(*info)
	if err != nil {
		network.ErrorShow(w, err)
		return
	}
	network.ResultOK(w, 0, "小说资源获取成功！", info)
}

func (p *resourcesController) getSourcesRequest(info entity.NovelInfoEntity) (*entity.NovelInfoEntity, error) {
	handler := bookcity.SelectResources(info.Source)
	html, err := handler.RequestSearchHTML(info.Title)
	if err != nil {
		return nil, err
	}
	info.SourceURL, err = handler.SourceURLHTMLParse(html, info.Title, info.Author)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (p *resourcesController) SourcesDBOperation(info entity.NovelInfoEntity) error {
	//如果不存在该小说数据，就插入；存在就更新
	if dao.GetInfoDaoInstance().HasNovelInfo(info.ID) {
		if id := dao.GetInfoDaoInstance().InsertNovelInfo(info); id <= 0 {
			return network.Errors{Code: -301, Msg: network.ErrorCode[-301]}
		} else {
			return nil
		}
	} else {
		if id := dao.GetInfoDaoInstance().UpdateNovelInfo(info); id <= 0 {
			return network.Errors{Code: -302, Msg: network.ErrorCode[-302]}
		} else {
			return nil
		}
	}
}
