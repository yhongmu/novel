package controller

import (
	"Novel/entity"
	"net/http"
	"net/url"

	"Novel/network"
)

const QIDIAN_SEARCH_URL = "https://www.qidian.com/search?kw="

var search *searchController

func GetSearchInstance() *searchController {
	if search == nil {
		search = &searchController{}
	}
	return search
}

type searchController struct {
}

func (p *searchController) Router(router *network.RouterHandler)  {
	router.Router("/novel/search", p.search)
}

func (p *searchController) search(w http.ResponseWriter, r *http.Request) {
	kw := r.PostFormValue("km")
	if kw == "" {
		network.ResultFail(w, -101, network.ErrorCode[-101])
		return
	}
	if searchEntity, err := p.searchRequest(kw); err != nil {
		network.ErrorShow(w, err)
	} else {
		network.ResultOK(w, 0, "搜索成功！", searchEntity)
	}
}

func (*searchController) searchRequest(kw string) ([]entity.NovelInfoEntity, error) {
	kw = url.QueryEscape(kw)
	searchURL := QIDIAN_SEARCH_URL + kw
	if searchEntity, err := network.GetSearchRequest(searchURL, network.SearchResultHTMLParse); err == nil {
		if search, ok := searchEntity.([]entity.NovelInfoEntity); ok {
			return search, nil
		} else {
			return nil, network.Errors{Code: 1, Msg: "cw"}
		}
	} else {
		return nil, err
	}
}
