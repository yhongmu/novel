package main

import (
	"Novel/controller"
	"Novel/dao"
	"Novel/network"
	"fmt"
	"net/http"
	"time"
)

func main(){
	dao.InitDB()
	dao.CreateNovelInfoTable()
	server := http.Server{
		Addr: ":8080",
		Handler: network.GetRouterInstance(),
		ReadTimeout: 5 * time.Second,
	}
	//注册搜索API路由
	controller.GetSearchInstance().Router(network.GetRouterInstance())
	//注册获取小说资源API路由
	controller.GetResourcesInstance().Router(network.GetRouterInstance())
	//注册获取小说目录API路由
	controller.GetDirectoryInstance().Router(network.GetRouterInstance())
	//注册获取小说内容API路由
	controller.GetContentInstance().Router(network.GetRouterInstance())

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("start server error")
	}
	fmt.Println("start server success")
}
