package util

import (
	"github.com/ooncn/common/oredis"
	"strings"
	"time"
)

type CodeCheck struct {
	Tag   string `json:"tag"`   //标签
	Code  string `json:"code"`  //验证码
	Del   int    `json:"del"`   //是否删除 0否、1是
	Phone string `json:"phone"` //手机号
	Md5   string `json:"md5"`
}

func (c *CodeCheck) GetMd5() string {
	return Md5Str(c.Md5)
}

/**
 * 生产验证码
 *
 * @param account    手机号/账号
 * @param codeLength 验证码长度
 * @param type       0.数字/1.英文/2.数字英文混合
 * @param validTime  保留有效时间（秒）
 * @return code验证码
 */
func (c *CodeCheck) SetCode(codeLength, codeType int) {
	var code string
	switch codeType {
	case 0:
		code = RandInt(codeLength)
		break
	case 1:
		code = RandAa(codeLength)
		break
	case 2:
		code = RandAaInt(codeLength)
		break
	}
	c.Code = strings.ToUpper(code)
}

/**
 * 生产验证码
 *
 * @param key    redisKey
 * @param validTime  保留有效时间（秒）
 * @return code验证码
 */
func (c *CodeCheck) RedisSetCode(key string, validTime int64) {
	oredis.Ser.SetType(key, c, time.Duration(validTime)*time.Second)
} //Code sndCode(String account, Integer codeLength, Integer type, Integer validTime)

/**
 * 获取验证码
 *
 * @param account 手机号/账号
 * @return Code
 */
func (c CodeCheck) RedisGetCode(key string) CodeCheck {
	var d CodeCheck
	oredis.Ser.GetType(key, &d)
	return d
} //Code getCode(String account)
func (c *CodeCheck) RedisDelCode(key string) {
	oredis.Ser.Del(key)
} //void delCode(String account)

/**
 * 检验验证码
 *
 * @param account 手机号/账号
 * @param code    验证码
 * @return Boolean false/true
 */
func (c *CodeCheck) VerificationCode(key, code string) bool {
	m := c.RedisGetCode(key)
	return m.Code == strings.ToUpper(code)

} //Boolean verificationCode(String account, String code)
