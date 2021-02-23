package util

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"
)

type EncryptUtil struct{}

/**
AES-CBC 加密
*/
func AesCBCNewEncrypt(orig string) string {
	return AesCBCNewByteEncrypt([]byte(orig))
}
func AesCBCNewByteEncrypt(orig []byte) string {
	return AesCBCEncryptByteKey(orig, "weixuan"+time.Now().Format("2006-01-02"))
}
func AesCBCEncryptByte(origData []byte, k []byte) string {
	a := &AES{}
	i, err := a.AesEncrypt(origData, k, k)
	if err != nil {
		return ""
	}
	return i
}
func AesCBCEncryptByteKey(origData []byte, key string) string {
	k := []byte(key)
	return AesCBCEncryptByte(origData, k)
}
func AesCBCEncrypt(orig string, key string) string {
	defer func() {
		a := recover()
		if a != nil {
			return
		}
	}()
	//fmt.Println(key)
	//fmt.Println(len(key))
	// 转成字节数组
	origData := []byte(orig)
	return AesCBCEncryptByteKey(origData, key)
}

/**
AES-CBC 解密
*/
func AesCBCNewDecrypt(cry string) string {
	return AesCBCDecrypt(cry, "weixuan"+time.Now().Format("2006-01-02"))
}
func AesCBCDecrypt(cry string, key string) string {
	defer func() {
		a := recover()
		if a != nil {
			fmt.Println(a)
		}
	}()
	i := AesCBCDecryptByte(cry, key)
	return string(i)
}

func AesCBCDecryptByte(cry string, key string) []byte {
	a := &AES{}
	k := []byte(key)
	i, err := a.AesDecrypt(cry, k, k)
	if err != nil {
		return nil
	}
	return i
}

/**
MD5 加密
*/
func Md5Byte(data []byte) string {
	h := md5.New()
	h.Write(data) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(cipherStr))
}

/**
MD5字符串 加密
*/
func Md5Str(data string) string {
	h := md5.New()
	h.Write([]byte(data)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(cipherStr))
}
func PassMd5(data string) string {
	h := md5.New()
	h.Write([]byte(data)) // 需要加密的字符串为 123456
	h.Write([]byte("weixuan"))
	cipherStr := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(cipherStr))
}

/**
MD5salt字符串 加密
*/
func Md5Salt(data, salt string) string {
	h := md5.New()
	h.Write([]byte(data)) // 需要加密的字符串为 123456
	h.Write([]byte(salt))
	cipherStr := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(cipherStr))
}

/**
获取文件MD5值
*/
func Md5File(data os.File) (m string, size int64) {
	h := md5.New()
	size, _ = io.Copy(h, &data)
	cipherStr := h.Sum(nil)
	m = strings.ToUpper(hex.EncodeToString(cipherStr))
	return
}

func Md5FilePath(path string) (m string, size int64) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0
	}
	defer f.Close()
	return Md5File(*f)
}

func PasswordMd5Salt(pass string) string {
	return Md5Salt(pass, "weixuan2019")
}

/*func CodeCaptchaCreate() {
	//config struct for digits
	//数字验证码配置
	var configD = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}
	//config struct for audio
	//声音验证码配置
	var configA = base64Captcha.ConfigAudio{
		CaptchaLen: 6,
		Language:   "zh",
	}
	//config struct for Character
	//字符,公式,验证码配置
	var configC = base64Captcha.ConfigCharacter{
		Height: 60,
		Width:  240,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeNumber,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         6,
	}
	//创建声音验证码
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	idKeyA, capA := base64Captcha.GenerateCaptcha("", configA)
	//以base64编码
	base64stringA := base64Captcha.CaptchaWriteToBase64Encoding(capA)
	//创建字符公式验证码.
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	idKeyC, capC := base64Captcha.GenerateCaptcha("", configC)
	//以base64编码
	base64stringC := base64Captcha.CaptchaWriteToBase64Encoding(capC)

	//创建数字验证码.
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	idKeyD, capD := base64Captcha.GenerateCaptcha("", configD)
	//以base64编码
	base64stringD := base64Captcha.CaptchaWriteToBase64Encoding(capD)

	fmt.Println(configD)

	fmt.Println(idKeyA, base64stringA, "\n")
	fmt.Println(idKeyC, base64stringC, "\n")
	fmt.Println(idKeyD, base64stringD, "\n")
}*/

func PuEncrypt(s string) string {
	return AesCBCEncrypt(s, "weixuan2019-0924")
}

func PuDecrypt(s string) string {
	return string(PuDecryptToByte(s))
}
func PuDecryptToByte(s string) []byte {
	return AesCBCDecryptByte(s, "weixuan2019-0924")
}

func PuEncryptByte(s []byte) string {
	return AesCBCEncryptByteKey(s, "weixuan2019-0924")
}

type GoogleAuth struct {
}

func NewGoogleAuth() *GoogleAuth {
	return &GoogleAuth{}
}

func (this *GoogleAuth) un() int64 {
	return time.Now().UnixNano() / 1000 / 30
}

func (this *GoogleAuth) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func (this *GoogleAuth) base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func (this *GoogleAuth) base32decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func (this *GoogleAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (this *GoogleAuth) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}

func (this *GoogleAuth) oneTimePassword(key []byte, data []byte) uint32 {
	hash := this.hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := this.toUint32(hashParts)
	return number % 1000000
}

// 获取秘钥
func (this *GoogleAuth) GetSecret() string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, this.un())
	return strings.ToUpper(this.base32encode(this.hmacSha1(buf.Bytes(), nil)))
}

// 获取动态码
func (this *GoogleAuth) GetCode(secret string) (string, error) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := this.base32decode(secretUpper)
	if err != nil {
		return "", err
	}
	number := this.oneTimePassword(secretKey, this.toBytes(time.Now().Unix()/30))
	return fmt.Sprintf("%06d", number), nil
}

// 获取动态码二维码内容
func (this *GoogleAuth) GetQrcode(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}

// 获取动态码二维码图片地址,这里是第三方二维码api
func (this *GoogleAuth) GetQrcodeUrl(user, secret string) string {
	qrcode := this.GetQrcode(user, secret)
	width := "200"
	height := "200"
	data := url.Values{}
	data.Set("data", qrcode)
	return "https://api.qrserver.com/v1/create-qr-code/?" + data.Encode() + "&size=" + width + "x" + height + "&ecc=M"
}

var err error

// 验证动态码
func (this *GoogleAuth) VerifyCode(secret, code string) (bool, error) {
	_code, err := this.GetCode(secret)
	if err != nil {
		return false, err
	}
	fmt.Println(_code, secret, code)
	return _code == code, nil
}

type Rsa2 struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// 生成密钥对
func CreateKeys(publicKeyWriter, privateKeyWriter io.Writer, keyLength int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return err
	}
	derStream := MarshalPKCS8PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	err = pem.Encode(privateKeyWriter, block)
	if err != nil {
		return err
	}

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	err = pem.Encode(publicKeyWriter, block)
	if err != nil {
		return err
	}

	return nil
}

func NewRsa2PKCS8(publicKey string, privateKey string) (*Rsa2, error) {
	publicKey = fmt.Sprintf(`-----BEGIN RSA PUBLIC KEY-----
%s
-----END RSA PUBLIC KEY-----`, publicKey)
	privateKey = fmt.Sprintf(`-----BEGIN RSA PRIVATE KEY-----
%s
-----END RSA PRIVATE KEY-----`, privateKey)

	return NewRsa2PKCS8Byte([]byte(publicKey), []byte(privateKey))
}
func NewRsa2PKCS8Byte(publicKey []byte, privateKey []byte) (*Rsa2, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	block, _ = pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pri, ok := priv.(*rsa.PrivateKey)
	if ok {
		return &Rsa2{
			publicKey:  pub,
			privateKey: pri,
		}, nil
	} else {
		return nil, errors.New("private key not supported")
	}
}
func NewRsa2PKCS1(publicKey string, privateKey string) (*Rsa2, error) {
	publicKey = fmt.Sprintf(`-----BEGIN RSA PUBLIC KEY-----
%s
-----END RSA PUBLIC KEY-----`, publicKey)
	privateKey = fmt.Sprintf(`-----BEGIN RSA PRIVATE KEY-----
%s
-----END RSA PRIVATE KEY-----`, privateKey)

	return NewRsa2PKCS1Byte([]byte(publicKey), []byte(privateKey))
}
func NewRsa2PKCS1Byte(publicKey []byte, privateKey []byte) (*Rsa2, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	block, _ = pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &Rsa2{
		publicKey:  pub,
		privateKey: priv,
	}, nil
}

// 公钥加密
func (r *Rsa2) PublicEncrypt(data string) (string, error) {
	partLen := r.publicKey.N.BitLen()/8 - 11
	chunks := split([]byte(data), partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		v15, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(v15)
	}

	return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil
}

// 私钥解密
func (r *Rsa2) PrivateDecrypt(encrypted string) (string, error) {
	partLen := r.publicKey.N.BitLen() / 8
	raw, err := base64.RawURLEncoding.DecodeString(encrypted)
	chunks := split(raw, partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(decrypted)
	}

	return buffer.String(), err
}

// 数据加签
func (r *Rsa2) Sign(data string) (string, error) {
	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sign), err
}

// 数据验签
func (r *Rsa2) Verify(data string, sign string) error {
	h := crypto.SHA256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	decodedSign, err := base64.RawURLEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(r.publicKey, crypto.SHA256, hashed, decodedSign)
}

func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)

	k, _ := asn1.Marshal(info)
	return k
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
