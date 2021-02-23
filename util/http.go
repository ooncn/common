package util

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cheggaaa/pb"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

type HttpUtil struct{}

//region post 上传文件 url 请求地址 extraParams 请求参数 file 文件路径 fileFieldName 文件键名
func POSTFile(url string, extraParams map[string]string, file, fileFieldName string) string {

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	fileName := file[strings.LastIndex(file, "/")+1:]
	fileWriter1, _ := bodyWriter.CreateFormFile(fileFieldName, fileName)
	file1, _ := os.Open(file)
	defer file1.Close()
	io.Copy(fileWriter1, file1)

	for key, value := range extraParams {
		_ = bodyWriter.WriteField(key, value)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, _ := http.Post(url, contentType, bodyBuffer)
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	returnStr := string(respBody)
	return returnStr
}

//region post请求
func POSTParams(url string, extraParams map[string]string) (returnStr string, err error) {

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	for key, value := range extraParams {
		err = bodyWriter.WriteField(key, value)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuffer)
	if err != nil {
	}

	err = resp.Body.Close()
	if err != nil {
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	returnStr = string(respBody)
	return returnStr, err
}

//endregion
//region get请求
func GETParams(url string, extraParams map[string]string) (string, error) {
	var param string
	for key, value := range extraParams {
		param += key + "=" + value
	}
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}
	returnStr := string(respBody)
	return returnStr, err
}

//endregion
//region get URl请求
func GETToUrl(url string) (string, error) {
	respBody, err := GETUrlToByte(url)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}
	returnStr := string(respBody)
	return returnStr, err
}
func GETUrlToByte(url string) (respBody []byte, err error) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	return respBody, err
}

//endregion

/**
生成浏览器MD5token值
*/
func UserAgentAndIpToToken(r *http.Request) (string, string, string) {
	userAgent := r.Header.Get("User-Agent")
	ip := GetIP(r)
	token := Md5Salt((userAgent + ip + TimeUtil.DateToyMdHms()), "o0ncn")
	return userAgent, ip, token
}
func GetIP(c *http.Request) string {
	ip := c.RemoteAddr
	switch ip {
	case "127.0.0.1", "localhost":
		ip = c.Header.Get("X-Forwarded-For")
		break
	}
	return ip
}

type ReturnVo struct {
	Code interface{} `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//获取token
func GetToken(r *http.Request) string {
	return r.Header.Get("Token")
}
func ReturnJson(w http.ResponseWriter, r *http.Request, code int, msg string, object interface{}) {
	re := make(map[string]interface{})
	re["code"] = code
	re["msg"] = msg
	if IsNoBlank(object) {
		re["data"] = object
	}

	w.Header().Add("Access-Control-Allow-Methods", r.Method)
	w.Header().Add("Access-Control-Max-Age", "1800")
	w.Header().Add("Allow", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, JsonToStr(re))
}
func ReturnJsonCM(w http.ResponseWriter, r *http.Request, code int, msg string) {
	re := make(map[string]interface{})
	re["code"] = code
	re["msg"] = msg
	w.Header().Add("Access-Control-Allow-Methods", r.Method)
	w.Header().Add("Access-Control-Max-Age", "1800")
	w.Header().Add("Allow", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, JsonToStr(re))
}

type HttpOk struct {
	IsHttps       bool              // 是否是https
	Url           string            //连接地址
	Method        string            //请求方法 GET，POST
	Params        interface{}       //发送参数
	Form          map[string]string //发送表单参数
	ContentType   string            //发送请求类型
	Header        map[string]string // 请求头部
	Response      []byte            //返回主体二进制
	ResponseBody  string
	ContentLength int64    //返回主体数据长度
	Advance       chan int //进度
	StatusCode    int
	TimeOut       time.Duration //超时时间
	ResponseObj   *http.Response
	File          *os.File
	FileInputName string
	FileName      string
	BodyBuffer    *bytes.Buffer
	ReqBodyMutex  sync.Mutex
	ResqBodyMutex sync.Mutex
}

func (h *HttpOk) Query() (err error) {
	var bytesData []byte
	if h.Params != nil {
		t := reflect.TypeOf(h.Params)
		v := reflect.ValueOf(h.Params)
		switch t.Kind() {
		case reflect.String:
			bytesData = []byte(v.String())
			break
		case reflect.Bool:
			if v.Bool() {
				bytesData = []byte("true")
			} else {
				bytesData = []byte("false")
			}
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var buffer bytes.Buffer
			err = binary.Write(&buffer, binary.BigEndian, v.Int())
			bytesData = buffer.Bytes()
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			var buffer bytes.Buffer
			err = binary.Write(&buffer, binary.BigEndian, v.Uint())
			bytesData = buffer.Bytes()
			break
		case reflect.Float32, reflect.Float64:
			var buffer bytes.Buffer
			err = binary.Write(&buffer, binary.BigEndian, v.Float())
			bytesData = buffer.Bytes()
			break
		default:
			bytesData, err = json.Marshal(h.Params)
		}
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest(h.Method, h.Url, reader)
	if err != nil {
		return err
	}
	for i := 0; i < len(h.Header); i++ {
		for key, value := range h.Header {
			request.Header.Set(key, value)
		}
	}
	request.Header.Set("Content-Type", h.ContentType)
	client := http.Client{}
	if h.IsHttps {
		client.Transport = &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		}
	}
	if h.TimeOut > 0 {
		client.Timeout = h.TimeOut
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("连接失败")
	}
	h.StatusCode = resp.StatusCode
	h.ResponseObj = resp
	var respBytes []byte
	contentLength := resp.ContentLength
	if contentLength == 0 {
		h.ResqBodyMutex.Lock()
		respBytes, err = ioutil.ReadAll(resp.Body)
		h.ResqBodyMutex.Unlock()
		if err != nil {
			return err
		}
	} else {
		h.ResqBodyMutex.Lock()
		h.ContentLength = contentLength
		body := bufio.NewReader(resp.Body)
		defer func() {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()
		//bar := pb.StartNew(int(h.ContentLength))
		for {
			content, err := body.ReadByte()
			if err == io.EOF {
				break
			}
			//bar.Increment()
			respBytes = append(respBytes, content)
		}
		//bar.Finish()
		h.ResqBodyMutex.Unlock()
	}

	h.Response = respBytes
	respStr := string(respBytes)
	if strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), "gbk") {
		respStr = GBKConvertUTF8(respStr, "gbk", "utf-8")
	}
	h.ResponseBody = respStr
	return err
}
func (h *HttpOk) QueryJson(jsonString string) (err error) {
	return h.QueryJsonByte([]byte(jsonString))
}
func (h *HttpOk) QueryBody(body string) (err error) {
	return h.QueryBodyByte([]byte(body))
}
func (h *HttpOk) QueryJsonByte(jsonByte []byte) (err error) {
	request, err := http.NewRequest("POST", h.Url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return err
	}
	for i := 0; i < len(h.Header); i++ {
		for key, value := range h.Header {
			request.Header.Set(key, value)
		}
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	if h.TimeOut > 0 {
		client.Timeout = h.TimeOut
	}
	resp, err := client.Do(request)

	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("连接失败")
	}
	h.StatusCode = resp.StatusCode
	h.ResponseObj = resp
	contentLength := resp.ContentLength
	if contentLength < 0 {
		contentLength = 1
	}
	h.ContentLength = contentLength

	body := bufio.NewReader(resp.Body)
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	dSize := contentLength / 1024 / 1024
	respBytes := make([]byte, contentLength)
	if dSize <= 0 {
		respBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	} else {
		i := 0
		bar := pb.StartNew(int(h.ContentLength))
		for {
			content, err := body.ReadByte()
			if err == io.EOF {
				break
			}
			bar.Increment()
			respBytes[i] = content
			i++
		}
		bar.Finish()
	}

	//byte数组直接转成string，优化内存
	if len(respBytes) < 1 {
		return
	}
	h.Response = respBytes
	h.ResponseBody = string(respBytes)
	return err
}
func (h *HttpOk) QueryBodyByte(bodyByte []byte) (err error) {
	request, err := http.NewRequest("POST", h.Url, bytes.NewBuffer(bodyByte))
	if err != nil {
		return err
	}
	for i := 0; i < len(h.Header); i++ {
		for key, value := range h.Header {
			request.Header.Set(key, value)
		}
	}
	client := http.Client{}
	if h.TimeOut > 0 {
		client.Timeout = h.TimeOut
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	h.StatusCode = resp.StatusCode
	h.ResponseObj = resp
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(respBytes) < 1 {
		return
	}
	//byte数组直接转成string，优化内存
	h.Response = respBytes
	h.ResponseBody = string(respBytes)
	return err
}
func (h *HttpOk) QueryGetTimeOut() (err error) {
	if h.Form != nil {
		data := make(url.Values)
		for k, v := range h.Form {
			data.Set(k, v)
		}
		h.Url = h.Url + "?" + data.Encode()
	}
	request, err := http.NewRequest("GET", h.Url, nil)
	if err != nil {
		return err
	}

	for i := 0; i < len(h.Header); i++ {
		for key, value := range h.Header {
			request.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	if h.TimeOut > 0 {
		client.Timeout = h.TimeOut
	}
	//client.Timeout = 10
	resp, err := client.Do(request)

	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("连接失败")
	}

	contentLength := resp.ContentLength
	if contentLength < 0 {
		contentLength = 1
	}
	h.ContentLength = contentLength

	body := bufio.NewReader(resp.Body)
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	dSize := resp.ContentLength / 1024 / 1024
	respBytes := make([]byte, contentLength)
	if dSize <= 0 {
		respBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	} else {
		i := 0

		bar := pb.StartNew(int(h.ContentLength))
		for {
			content, err := body.ReadByte()
			if err == io.EOF {
				break
			}
			bar.Increment()
			respBytes[i] = content
			i++
		}
		bar.Finish()
	}

	//byte数组直接转成string，优化内存
	if len(respBytes) < 1 {
		return
	}
	h.Response = respBytes
	h.ResponseBody = string(respBytes)
	return err
}
func (h *HttpOk) QueryGet() (err error) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	var src = h.Url
	if h.Form != nil {
		data := make(url.Values)
		for k, v := range h.Form {
			data.Set(k, v)
		}
		src = h.Url + "?" + data.Encode()
	}
	request, err := http.NewRequest("GET", src, nil)
	for i := 0; i < len(h.Header); i++ {
		for key, value := range h.Header {
			request.Header.Set(key, value)
		}
	}
	t := 5 * time.Second
	if h.TimeOut > t {
		t = h.TimeOut
	}
	client := &http.Client{Timeout: t}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	h.StatusCode = resp.StatusCode
	h.ResponseObj = resp
	//byte数组直接转成string，优化内存
	body := bufio.NewReader(resp.Body)
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	var respBytes = make([]byte, 0)
	if resp.ContentLength > 0 {
		dSize := resp.ContentLength / 1024 / 1024
		if dSize > 0 {
			bar := pb.StartNew(int(h.ContentLength))
			for {
				content, err := body.ReadByte()
				if err == io.EOF {
					break
				}
				bar.Increment()
				respBytes = append(respBytes, content)
			}
			bar.Finish()
		}
	} else {
		respBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	if len(respBytes) < 1 {
		return
	}
	//byte数组直接转成string，优化内存
	h.Response = respBytes
	h.ResponseBody = string(respBytes)
	return err
}
func (h *HttpOk) Get() (err error) {
	h.ResqBodyMutex.Lock()
	defer h.ResqBodyMutex.Unlock()
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	src, err := url.Parse(h.Url)
	if err != nil {
		return err
	}
	if h.Form != nil {
		data := src.Query()
		for k, v := range h.Form {
			data.Set(k, v)
		}
		src.RawQuery = data.Encode()
	}
	request, err := http.NewRequest("GET", src.String(), nil)
	for i := 0; i < len(h.Header); i++ {
		for key, value := range h.Header {
			request.Header.Set(key, value)
		}
	}
	t := 5 * time.Second
	if h.TimeOut > t {
		t = h.TimeOut
	}
	client := &http.Client{Timeout: t}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	h.StatusCode = resp.StatusCode
	h.ResponseObj = resp
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(respBytes) < 1 {
		return
	}
	//byte数组直接转成string，优化内存
	h.Response = respBytes
	h.ResponseBody = string(respBytes)
	return err
}
func (h *HttpOk) QueryGetData() (err error) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	h.Method = "GET"
	var src = h.Url
	if h.Form != nil {
		data := make(url.Values)
		for k, v := range h.Form {
			data.Set(k, v)
		}
		src = h.Url + "?" + data.Encode()
	}
	client := &http.Client{Timeout: 5 * time.Second}
	if h.TimeOut > 0 {
		client = &http.Client{Timeout: h.TimeOut * time.Second}
	}
	resp, err := client.Get(src)
	if err != nil {
		return err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	h.Response = respBytes
	h.ResponseBody = string(respBytes)
	return err
}
func (h *HttpOk) QueryForm() (err error) {
	data := make(url.Values)
	for k, v := range h.Form {
		data.Set(k, v)
	}
	resp, err := http.PostForm(h.Url, data)
	if err != nil {
		return err
	}
	//byte数组直接转成string，优化内存

	body := bufio.NewReader(resp.Body)
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	dSize := resp.ContentLength / 1024 / 1024
	respBytes := make([]byte, resp.ContentLength)
	if dSize <= 0 {
		respBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	} else {
		i := 0

		bar := pb.StartNew(int(h.ContentLength))
		for {
			content, err := body.ReadByte()
			if err == io.EOF {
				break
			}
			bar.Increment()
			respBytes[i] = content
			i++
		}
		bar.Finish()
	}
	//byte数组直接转成string，优化内存
	if len(respBytes) < 1 {
		return
	}
	h.Response = respBytes
	h.ResponseBody = string(respBytes)
	return err
}
func (h *HttpOk) PostForm() (err error) {
	data := make(url.Values)
	for k, v := range h.Form {
		data.Set(k, v)
	}
	resp, err := http.PostForm(h.Url, data)
	if err != nil {
		return err
	}
	for i := 0; i < len(h.Header); i++ {
		for key, value := range h.Header {
			resp.Header.Set(key, value)
		}
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//byte数组直接转成string，优化内存
	if len(respBytes) < 1 {
		return
	}
	h.Response = respBytes
	h.ResponseBody = string(respBytes)
	return err
}
func (h *HttpOk) HeadQuery() (err error) {
	request, err := http.NewRequest(h.Method, h.Url, nil)
	for i := 0; i < len(h.Header); i++ {
		for key, value := range h.Header {
			request.Header.Set(key, value)
		}
	}
	client := http.Client{}
	if h.IsHttps {
		client.Transport = &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		}
	}
	if h.TimeOut > 0 {
		client.Timeout = h.TimeOut
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//byte数组直接转成string，优化内存
	if len(respBytes) < 1 {
		return
	}
	h.StatusCode = resp.StatusCode
	h.ResponseObj = resp
	h.Response = respBytes
	h.ResponseBody = string(respBytes)
	return err
}
func QueryFormFileBodyBuffer(path, FileInputName, FileName string) (contentType string, bodyBuffer *bytes.Buffer, err error) {
	bodyBuffer = &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	defer bodyWriter.Close()
	fileWriter, _ := bodyWriter.CreateFormFile(FileInputName, FileName)
	file2, _ := os.Open(path)
	defer file2.Close()
	io.Copy(fileWriter, file2)
	contentType = bodyWriter.FormDataContentType()
	return
}
func (h *HttpOk) QueryFormFile(contentType string, bodyBuffer *bytes.Buffer) (err error) {
	request, err := http.NewRequest("POST", h.Url, bodyBuffer)
	for i := 0; i < len(h.Header); i++ {
		for key, value := range h.Header {
			request.Header.Set(key, value)
		}
	}
	request.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBytes, _ := ioutil.ReadAll(resp.Body)
	//byte数组直接转成string，优化内存
	if len(respBytes) < 1 {
		return
	}
	h.StatusCode = resp.StatusCode
	h.ResponseObj = resp
	h.Response = respBytes
	h.ResponseBody = string(respBytes)
	return err
}

// TODO html主体图片转换base64
func HtmlBodyImgToBase64(str string) string {
	imgs := strings.Split(str, "&lt;img ")
	if len(imgs) < 2 {
		return str
	}
	var content string
	for _, v := range imgs {
		src := "src=&#34;"
		i := strings.Index(v, src)
		// 判断是否符合图片请求
		if i >= 0 {
			content += "&lt;img " + v[:i] + src
			i += len(src)
			//证明该字符串中存在图片文件
			jc := v[i:]
			j := strings.Index(jc, "&#34;")
			if j >= 0 {
				imgUrl := jc[:j]
				if strings.Index(imgUrl, `//`) == 0 {
					imgUrl = "http:" + imgUrl
				}
				imgBase64, err := HttpGetImg(imgUrl)
				if err != nil {
					content += imgUrl + jc[j:]
				} else {
					content += imgBase64 + jc[j:]
				}
				//img(f)
			} else {
				content += jc
			}
		} else {
			content += v
		}
	}
	return content
}

// TODO html主体图片转换base64
func HtmlBodyImgToFile(str, filePath string, hostArr []string) string {
	imgs := strings.Split(str, "&lt;img ")
	if len(imgs) < 2 {
		return str
	}
	var content string
	for _, v := range imgs {
		src := "src=&#34;"
		i := strings.Index(v, src)
		// 判断是否符合图片请求
		if i >= 0 {
			content += "&lt;img " + v[:i] + src
			i += len(src)
			//证明该字符串中存在图片文件
			jc := v[i:]
			j := strings.Index(jc, "&#34;")
			if j >= 0 {
				imgUrl := jc[:j]
				if strings.Index(imgUrl, `//`) == 0 {
					imgUrl = "http:" + imgUrl
				}
				u, err := url.Parse(imgUrl)
				if err == nil {
					if !ArrContainsStr(hostArr, u.Host) {
						imgBase64, err := HttpGetImg(imgUrl)
						if err == nil {
							content += imgBase64 + jc[j:]
						} else {
							content += imgUrl + jc[j:]
						}
					} else {
						content += imgUrl + jc[j:]
					}
				} else {
					content += imgUrl + jc[j:]
				}
			} else {
				content += jc
			}
		} else {
			content += v
		}
	}
	return content
}

func ArrContainsStr(arr []string, v string) bool {
	if arr != nil && len(arr) > 0 {
		for _, k := range arr {
			if k == v {
				return true
			}
		}
	}
	return false
} // 字符串数组是否有相同的字符串

// TODO 根据图片路径获取图片base64值
func HttpGetImg(imgUrl string) (imgBase64 string, err error) {
	h := HttpOk{Url: imgUrl}
	err = h.Get()
	if err != nil {
		return
	}
	resp := h.ResponseObj
	hander := resp.Header
	fileType := hander["Content-Type"][0]
	if len(fileType) < 6 || strings.Index(fileType, "image") < 0 {
		fileType = "image/jpeg"
	}
	f := h.Response
	imgBase64 = base64.StdEncoding.EncodeToString(f)
	imgBase64 = "data:" + fileType + ";base64," + imgBase64
	return
}
func HttpGetImgToFile(filePath, imgUrl string) (err error) {
	h := HttpOk{Url: imgUrl}
	err = h.Get()
	if err != nil {
		return
	}
	resp := h.ResponseObj
	header := resp.Header
	fileType := header["Content-Type"][0]
	if len(fileType) < 6 || strings.Index(fileType, "image") < 0 {
		fileType = "image/jpeg"
	}
	f := h.Response
	fr, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer fr.Close()
	if err != nil {
		return
	} else {
		_, err = fr.Write(f)
	}
	return
}

func UrlParamMap(m map[string]interface{}) (s string) {
	if m == nil {
		return
	}
	for k, v := range m {
		s += "&" + k + "=" + JsonToStr(v)
	}
	return s[1:]
}
func UrlAddParam(url string, m map[string]interface{}) string {
	if strings.Index(url, "?") > 1 {
		return url + "&" + UrlParamMap(m)
	} else {
		return url + "?" + UrlParamMap(m)
	}

}
func AddCertConfig(certFilePath, keyFilePath, pkcs12FilePath interface{}) (tlsConfig *tls.Config, err error) {
	if certFilePath == nil && keyFilePath == nil && pkcs12FilePath == nil {
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
		return tlsConfig, nil
	}

	if certFilePath != nil && keyFilePath != nil && pkcs12FilePath != nil {
		cert, err := ioutil.ReadFile(certFilePath.(string))
		if err != nil {
			return nil, fmt.Errorf("ioutil.ReadFile：%w", err)
		}
		key, err := ioutil.ReadFile(keyFilePath.(string))
		if err != nil {
			return nil, fmt.Errorf("ioutil.ReadFile：%w", err)
		}
		pkcs, err := ioutil.ReadFile(pkcs12FilePath.(string))
		if err != nil {
			return nil, fmt.Errorf("ioutil.ReadFile：%w", err)
		}
		pkcsPool := x509.NewCertPool()
		pkcsPool.AppendCertsFromPEM(pkcs)
		certificate, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, fmt.Errorf("tls.LoadX509KeyPair：%w", err)
		}
		tlsConfig = &tls.Config{
			Certificates:       []tls.Certificate{certificate},
			RootCAs:            pkcsPool,
			InsecureSkipVerify: true}
		return tlsConfig, nil
	}
	return nil, errors.New("cert paths must all nil or all not nil")
}
