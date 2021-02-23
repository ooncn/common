package util

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

/**
检测环境
*/

type PcReq struct {
	Url    string `json:"url"`
	Id     string `json:"id"`
	IdCard string `json:"idCard"`
	Reader int    `json:"reader"`
}

type PcResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func JsoniterToType(config string, result interface{}) (err error) {
	api := jsoniter.ConfigCompatibleWithStandardLibrary
	err = api.Unmarshal([]byte(config), result)
	return
}

/*
a 交易金额
b 计算费率
y 剩余金额
c 手续费
*/
func CommissionPrice(a, b float64) (y, c float64) {
	fmt.Println("交易金额", a)
	fmt.Println("计算费率", b/100)
	str := fmt.Sprintf("%.2f", a*b/100)
	fmt.Println("手续费")
	c, _ = strconv.ParseFloat(str, 64)
	return a - c, c
}

var (
	EventCode = map[int]string{
		0: "授权",
		1: "发送消息",
		2: "在线人数",
		3: "系统通知",
	}
	ToCode = map[int]string{
		0: "系统",
		//1: "发送消息",
		//2: "在线人数",
		//3: "系统通知",
	}
)

const (
	WS_EVENT_CODE_TOKEN  = iota // 授权
	WS_EVENT_CODE_MSG           // 发送消息
	WS_EVENT_CODE_ONLINE        // 在线人数
	WS_EVENT_CODE_SYS           // 系统通知
)

type WsReq struct {
	Token string `json:"token"` // 授权token
	To    string `json:"to"`    // 接受消息的人
	Event int    `json:"event"` // 事件：0、授权 1、发送消息
	Data  string `json:"data"`  // 内容
	//Route string `json:"reader"` // 路由
}
