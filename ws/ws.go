package ws

import (
	"github.com/ooncn/common/util"
)

type WsObj struct {
	Type  string      `json:"type"`  //消息类型
	Token string      `json:"token"` //用户名:
	Code  int         `json:"code"`  //消息编码 0正常，1成功，2警告，3错误
	Data  interface{} `json:"data"`  //消息主体
}

/**
 * 发送前端消息
 * @param {[type]} t     [消息类型]
 * @param {[type]} token string        [用户名]
 * @param {[type]} code  int           [编码]
 * @param {[type]} data  interface{} [消息主体]
 */
func NewWsObj(logType, token string, code int, data interface{}) {
	Publish <- NewEvent(EVENT_MESSAGE, token, []byte(util.JsonToStr(WsObj{logType, token, code, data})))
}
