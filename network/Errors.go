package network

import "fmt"

type Errors struct {
	Code int16
	Msg string
}

func (e Errors) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

var ErrorCode = map[int16]string{
	0: "成功",
	-1: "表示系统错误",

	-100: "传参错误",
	-101: "搜索框未输入信息！",
	-102: "未输入小说id或title或author！",

	-200: "html解析错误",
	-201: "未在选中书城搜索到该小说",
	-202: "未获取到该小说目录",

	-300: "数据库操作错误",
	-301: "数据插入错误",
	-302: "数据更新错误",
	-303: "未查询到章节URL",

	-400: "网络请求错误",
	-401: "同步网络请求错误",
}
