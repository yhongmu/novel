package network

import "net/http"

var router *RouterHandler

func GetRouterInstance() *RouterHandler {
	if router == nil {
		router = &RouterHandler{
			mux: make(map[string]func (http.ResponseWriter, *http.Request)),
		}
	}
	return router
}

type RouterHandler struct {
	mux map[string]func (http.ResponseWriter, *http.Request)
}

func (p *RouterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if fun, ok := p.mux[r.URL.Path]; ok {
		fun(w, r)
		return
	}
	http.Error(w, "error URL:" + r.URL.String(), http.StatusBadRequest)
}

func (p *RouterHandler) Router(relativePath string, handler func(w http.ResponseWriter, r *http.Request)) {
	p.mux[relativePath] = handler
}