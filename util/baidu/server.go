package db

import (
	"encoding/base64"
	"fmt"
	"github.com/ooncn/common/util"
	"io/ioutil"
	"strings"
	"time"
)

/*
func BaiduVOP(speech string) {
	uri := "http://vop.baidu.com/server_api"
	{
		"format":"wav",
		"rate":16000,
		"dev_pid":1537,
		"channel":1,
		"token":xxx,
		"cuid":"baidu_workshop",
		"len":4096,
		"speech":"xxx", // xxx为 base64（FILE_CONTENT）
	}
}*/

/*

   private static final String APP_ID = "14365274";
   private static final String API_KEY = "tnt2Pn8uwsG46xLMZk58gNa6";
   private static final String SECRET_KEY = "nhHb0BCvIZ9sGfrhjiyutoxjREPXdBBi";
*/
type BaiduOauth struct {
	AppID     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	SecretKey string `json:"secret_key"`
	BaiduAuth
}
type BaiduAuth struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (b *BaiduOauth) GetToken() {
	var auth BaiduAuth
	uri := "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=" +
		b.AppKey + "&client_secret=" +
		b.SecretKey
	get := util.HttpOk{Url: uri}
	get.Get()
	util.JsonToType(get.ResponseBody, &auth)
	b.BaiduAuth = auth
}

func (b *BaiduOauth) TTS(path string) {
	var (
		auth   = b.BaiduAuth
		speech string
	)
	w := path[strings.LastIndex(path, "/")+1:] // 获取最后
	fmt.Println(path)
	fmt.Println(w)
	i := strings.LastIndex(w, ".")
	fileNameExt := w[i+1:] // 获取最后
	fileName := w[:i]      // 获取最后
	fmt.Println(fileNameExt)
	fmt.Println(fileName)

	speech, _ = util.FileToBase64(path)

	uri := "http://vop.baidu.com/server_api"
	get := util.HttpOk{Url: uri, TimeOut: 10 * time.Minute}
	jsonStr := fmt.Sprintf(`{
    "format":"m4a",
    "rate":16000,
    "dev_pid":1537,
    "channel":1,
    "cuid":"%s",
    "token":"%s",
    "speech":"%s"
}`, "weixuan", auth.AccessToken, speech)
	err := get.QueryJson(jsonStr)
	fmt.Println(err)

	fmt.Println(get.ResponseBody)
}
func FileToBase64(filepath string) (encodeString string, err error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	encodeString = base64.StdEncoding.EncodeToString(data)
	return
}
