package idValidator

import (
	"errors"
	"github.com/ooncn/common/util"
	"math"
	"strconv"
	"strings"
)

type IdCard struct {
	Sex     int    `json:"sex"`     //性别
	Length  int    `json:"length"`  //长度
	Address string `json:"address"` // 地址
	Zodiac  string `json:"zodiac"`  // 地址

	Body         string `json:"body"`         // 身份证号 body 部分
	AddressCode  string `json:"addressCode"`  //地址吗
	BirthdayCode string `json:"birthdayCode"` //出生日期码
	Order        string `json:"order"`        //出生日期码
	CheckBit     string `json:"checkBit"`     //校验码
	Type         string `json:"type"`         //出生日期码

	Year     int    `json:"year"`     // 年
	Province string `json:"province"` // 省份
	City     string `json:"city"`     // 城市
	District string `json:"district"` // 县区
}

func (i *IdCard) IsValidator(id string) (err error) {
	// 基础验证
	err = i.checkIdArgument(id)
	if err != nil {
		return errors.New("身份证校验错误")
	}
	// 分别验证：*地址码*、*出生日期码*和*顺序码*
	err = i.getAddressInfo()
	if err != nil {
		return err
	}
	err = i.checkBirthdayCode()
	if err != nil {
		return err
	}
	is := i.checkOrderCode()
	if !is {
		return errors.New("校验失败")
	}
	// 15位身份证不含校验码
	if i.Type == "15" {
		return nil
	}
	// 验证：校验码
	checkBit := i.generatorCheckBit()
	// 检查校验码
	if checkBit != i.CheckBit {
		err = errors.New("校验码错误")
	}
	return
}

func (i *IdCard) GetInfo(id string) (map[string]interface{}, error) {
	err := i.IsValidator(id)
	if err != nil {
		return nil, err
	}
	i.getZodiac()
	return map[string]interface{}{
		"addressCode":   i.AddressCode,
		"address":       i.Address,
		"birthdayCode":  i.BirthdayCode,
		"chineseZodiac": i.Zodiac,
		"sex":           i.Sex,
		"length":        i.Type,
		"checkBit":      i.CheckBit,
	}, nil
}

/*
检查并拆分身份证号
id 身份证号
*/
func (i *IdCard) checkIdArgument(id string) (err error) {
	// 将所有字符转大写
	id = strings.ToUpper(id)
	length := len(id) // 计算身份证号码长度
	if length == 15 {
		i.generateShortType(id)
		return
	} else if length == 18 {
		i.generateLongType(id)
		return
	}

	orderNum, _ := strconv.ParseInt(i.Order, 10, 64)
	if orderNum%2 == 0 {
		i.Sex = 0 //男
	} else {
		i.Sex = 1 //女
	}
	return errors.New("不合格")
}

func (i *IdCard) generateShortType(id string) {
	i.Body = id
	i.AddressCode = id[:6]
	i.BirthdayCode = "19" + id[6:12]
	year, _ := strconv.ParseInt("19"+id[6:8], 10, 64)
	i.Year = int(year)
	order := id[12:]
	i.Order = order
	i.Type = "15"
}
func (i *IdCard) generateLongType(id string) {
	i.Body = id[:17]
	i.AddressCode = id[:6]
	i.BirthdayCode = id[6:14]
	year, _ := strconv.ParseInt(id[6:10], 10, 64)
	i.Year = int(year)
	i.Order = id[14:17]
	i.CheckBit = id[17:]
	i.Type = "18"
}

//地址码
func (i *IdCard) getAddressInfo() (err error) {
	addressCode := i.AddressCode
	birthdayCode := i.BirthdayCode
	// 省级信息
	code := addressCode[:2] + "0000"
	address, is := i.getAddress(code, birthdayCode)
	var k = 0
	if is {
		k += 1
		i.Province = address
		i.Address += address
	}
	if k == 0 {
		return errors.New("地址码校验失败")
	}
	// 港澳台居民居住证无市级、县级信息
	if addressCode[:1] == "8" {
		return
	}
	// 市级信息
	code = addressCode[:4] + "00"
	address, is = i.getAddress(code, birthdayCode)
	if is {
		k += 1
		i.City = address
		i.Address += address
	}
	// 县级信息
	address, is = i.getAddress(addressCode, birthdayCode)
	if is {
		k += 1
		i.District = address
		i.Address += address
	}
	if k == 0 {
		return errors.New("地址码校验失败")
	}
	return
}

func (i *IdCard) getAddress(addressCode, birthdayCode string) (address string, is bool) {
	timeline, is := AddressCodeTimeline[addressCode]
	if is {
		year, _ := strconv.ParseInt(birthdayCode[:4], 10, 64)
		for k, v := range timeline {
			startYear, _ := strconv.ParseInt(v["startYear"], 10, 64)
			if k == "0" || year < startYear || year >= startYear {
				address, is = v["address"]
			}
		}
	}
	return
}

//出生日期码
func (i *IdCard) checkBirthdayCode() error {
	_, err := util.TimeUtil.FormatStr(i.BirthdayCode)
	return err
}

//顺序码
func (i *IdCard) checkOrderCode() bool {
	return len(i.Order) == 3
}

/*
 * 生成校验码
 * 详细计算方法 @lint https://zh.wikipedia.org/wiki/中华人民共和国公民身份号码
 * 身份证号 body 部分
 */
func (i *IdCard) generatorCheckBit() string {
	body := i.Body
	// 位置加权
	posWeight := make(map[int]int64)
	for z := 18; z > 1; z-- {
		posWeight[z] = int64(math.Pow(2, float64(z-1))) % 11
	}

	// 累身份证号 body 部分与位置加权的积
	var bodySum int64
	for k := range body {
		iu, _ := strconv.ParseInt(body[k:k+1], 10, 64)
		bodySum += iu * posWeight[18-k]
	}

	checkBit := (12 - (bodySum % 11)) % 11
	if checkBit == 10 {
		return "X"
	}
	return strconv.FormatInt(checkBit, 10)
}

/*
获取生肖信息
birthdayCode 出生日期码
*/
func (i *IdCard) getZodiac() {
	i.Zodiac = Zodiac[(i.Year-1900)%12]
}
