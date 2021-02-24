package network

import (
	"encoding/json"
	"net/http"
)

func ErrorShow(w http.ResponseWriter, err error) {
	if errDe, ok := err.(Errors); ok {
		//针对自定义的错误
		ResultFail(w, errDe.Code, errDe.Error())
	} else {
		//针对系统的错误
		ResultFail(w, -1, err.Error())
	}
}

func ResultOK(w http.ResponseWriter, Code int16, message string, data interface{}) {
	resultData := &ResultData{
		Code: Code,
		Message: message,
		Data: data,
	}
	ResultJson(w, resultData)
}

func ResultFail(w http.ResponseWriter, errorCode int16, errorMsg string) {
	resultData := &ResultData{
		Code: errorCode,
		Message: errorMsg,
		Data: nil,
	}
	ResultJson(w, resultData)
}

func ResultJson(w http.ResponseWriter, resultData *ResultData) {
	w.Header().Set("Content-Type", "application/json")
	jsonByte, _ := json.Marshal(resultData)
	_, _ = w.Write(jsonByte)
}

type ResultData struct {
	Code int16			`json:"code"`
	Message string		`json:"message"`
	Data interface{}	`json:"data,omitempty"`
}
