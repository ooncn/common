package util

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

func TestMd5Str(t *testing.T) {
	fmt.Println(Md5Str("1236456"))
	//A01A95CA9C0E2824AD9B07FD40916510
	fmt.Println(Md5Str("qwer.1324"))
	//6726BEAE503D84DD1F00919AC8A263F2
	fmt.Println(PassMd5("1236456"))
	//1A6017CA0FB06F649F7AE9DC0C85CB5C
	fmt.Println(PassMd5("qwer.1324"))
	//F93C4B2809089CA57DF74883E503A28C
}

var Pubkey = `-----BEGIN 公钥-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwjXTiZi7AJKDNfFRw5CjTj9VuHAZeCoo72nknCF2b4Aq6kh5KzghvTlPdh6sEYg2T7WaMEoS0EBhWZxqbKWSg1fYdrJtm0ZtHWhrsUJKG5kRw3D1JKgkPXWJznYqZsY0UdhnB1KJV0+dz7vEFTY/CTYRqULZua0DU5HqqA9CrfTt6sSYizbuwmPxVcLPURp2HkFfjp4mOy0TB6hKKbdnUOSJOZmoOmAOB/bGHRW2AsfHVlWbXuId/KI1WZBaGp1RbpoqgclGSRFBID+988F+qXWCsffrrOIypphnJdXYa3mIlurfaM+QvUCUEgr1UnI5g/+RxOMuZcv00/01+9/z4QIDAQAB
-----END 公钥-----
`

var Pirvatekey = `-----BEGIN 私钥-----
MIIEogIBAAKCAQEAwjXTiZi7AJKDNfFRw5CjTj9VuHAZeCoo72nknCF2b4Aq6kh5KzghvTlPdh6sEYg2T7WaMEoS0EBhWZxqbKWSg1fYdrJtm0ZtHWhrsUJKG5kRw3D1JKgkPXWJznYqZsY0UdhnB1KJV0+dz7vEFTY/CTYRqULZua0DU5HqqA9CrfTt6sSYizbuwmPxVcLPURp2HkFfjp4mOy0TB6hKKbdnUOSJOZmoOmAOB/bGHRW2AsfHVlWbXuId/KI1WZBaGp1RbpoqgclGSRFBID+988F+qXWCsffrrOIypphnJdXYa3mIlurfaM+QvUCUEgr1UnI5g/+RxOMuZcv00/01+9/z4QIDAQABAoIBAAX7qIetp90tplMsGwu+UfmgI+Dpuy8jhV1S0tMoiMCIn0nWb70wHoH42QTXjw9/NcUg4B4qugemZBlhozmzpB7dvvJxLsVA1y23wNGWLDhLb+uoeDTn5S3riBJPO6Es7AG1e/8SwN5qun7i1vicUjbTbpnbUF/S/648aZFB3xkCy9TzEAUAlD/s5bh1UI4necc30EnS716kBCDQoAlNNQdiRT+uEYay3nX7ySesMxLxgZURJ9tnW6+sIja/vQJSmc8Iy4x9sEK5GvSJaGjQnmHiVEuw3zBxAkHnQdn65KRxcwKPVvDa3jczsmCXHoG1D/5FmivY8ZW6l5jffkv7ds0CgYEA/LB/18WPJPzckwHwn64xS/wjJtouj0Mux1PGeqziBcKD3Z0byF2Tn6OeZ/0VutDEQElu1VYKwpwYxFLmsfmju3sT6TsH7c7bxdDDFGopBHtYZuu81JDo+i5c3QSb03JGLLo55wjkqctN+z45DPfw8DSN8XO94bnSH1AkwqEYh98CgYEAxMExOv6jMQvBO2nhBPlBK1AyUVleUiwvJCn1vHpESQJ7+luIPtyfJR92x/Dp83FWgwwMD6G6sJRrR1NA092RuT8lDkSySdf9u4OuPwBFyBh33wlpQk6KL5roxfJ4tEn1K+6m488kNGiDMneTo4aOw8CYv/V1QKrwaOtMkt3J/D8CgYAqlv0VOyEjVNNAm+UYpN1+NyMdm0yZrPMneYFMj/MQkXZ0VdSm8s6863D5ifitoh5Rz460umnZ30F1ZZuoh7EHGnmCqAZwGJuGPeeDe1kqfjeqMTWEhmAeOs5AGlTBUNNvGnxD6oXP8IpWPGiVPP0JH3KFLcLlVtKJoJJxk4F09QKBgGs/T2Fz6WpTmPmUxhYa853zceoLx7EM6olQ4eTh1JTjaMbX29VAFvN6ShnERRHwppJ6H5zpsESOMkfHpp+Vt9f9BmrXoUNFG8Z5iaJHuMHQLI8Dpz+AZix6yQUVHRxQ7/YJeSjWAUsb6N+6dFx/fRRQyDJiTo54XuEh2TR49p+rAoGAG+cr3lj8HpCJuvKXLOvMdsexcnWpvdAmRDMPU821kdSjr/OUjoIO+HbcBSsCIbIXKbWnn3QvERVOTPD/GiEBZxZnO0DhaPZvXbDYyBS7w6hDvthUsYaUJ8HcYmME7gjEztSVRI2E+bJh15A/GTa8mVA+t+OSEJH7hWexF/bsQfc=
-----END 私钥-----
`

func TestRSA2Test(t *testing.T) {
	// 公钥加密私钥解密
	if err := applyPubEPriD(); err != nil {
		log.Println(err)
	}
	// 公钥解密私钥加密
	if err := applyPriEPubD(); err != nil {
		log.Println(err)
	}
	fmt.Println()
	fmt.Println(gorsa.SignSha256WithRsa(
		`app_id=2021001168602207&biz_content={"scopes":["auth_base"],"state":"5bSRnzW9RJ"}&charset=UTF-8&format=json&sign_type=RSA2&timestamp=2020-06-12 14:48:55&version=1.0`,
		Pirvatekey))
}

// 公钥加密私钥解密
func applyPubEPriD() error {
	pubenctypt, err := gorsa.PublicEncrypt(`hello world`, Pubkey)
	if err != nil {
		return err
	}
	fmt.Printf("%s is:", pubenctypt)
	pridecrypt, err := gorsa.PriKeyDecrypt(pubenctypt, Pirvatekey)
	if err != nil {
		return err
	}
	fmt.Printf("%s is:", pridecrypt)
	if pridecrypt != `hello world` {
		return errors.New(`解密失败`)
	}
	return nil
}

// 公钥解密私钥加密
func applyPriEPubD() error {
	prienctypt, err := gorsa.PriKeyEncrypt(`hello world`, Pirvatekey)
	if err != nil {
		return err
	}

	pubdecrypt, err := gorsa.PublicDecrypt(prienctypt, Pubkey)
	if err != nil {
		return err
	}
	if string(pubdecrypt) != `hello world` {
		return errors.New(`解密失败`)
	}
	return nil
}
