package constant

import "fmt"

const (
	UiInvalidTimer                 = 172800 // 防伪码生成有效期（1天/86400秒)
	UiCreateNumMax                 = 100000 // 防伪码生成最大（1天/86400秒)
	UiCreateNumMin                 = 1000   // 防伪码生成最小（1天/86400秒)
	RedisKeyWxSnsAccessToken       = ":wx:sns:access_token"
	RedisKeyWxJsApiGetTicket       = ":wx:mq:JsApiGetTicket"
	RedisKeyWxMqOauth2AccessToken  = ":wx:mq:oauth2:access_token"
	REDIS_WX_USER_INFO             = "wx:userInfo:"
	RedisImgCode                   = "common:imgCode:"
	RedisWxMpConfig                = "system:wx:mp:config"
	REDIS_USER_ID                  = "user:id:"
	REDIS_USER_GID                 = "user:gid:"
	REDIS_USER_TYPE                = "user:type:"
	REDIS_USER_AUTH                = "user:auth:"
	REDIS_USER_LOGIN_LOG           = "user:login:log:"
	REDIS_USER_TOKEN               = "user:token:"
	REDIS_USER_SESSION             = "user:session:"
	REDIS_AGENT_SESSION            = "agent:session:"
	REDIS_MANAGER_USER_ID          = "manager:user:id:"
	REDIS_MANAGER_USER_GID         = "manager:user:gid:"
	REDIS_MANAGER_USER_TYPE        = "manager:user:type:"
	REDIS_MANAGER_USER_AUTH        = "manager:user:auth:"
	REDIS_MANAGER_USER_URLAUTH     = "manager:user:urlAuth:"
	REDIS_MANAGER_USER_MENU        = "manager:user:menu:"
	REDIS_MANAGER_USER_URLAUTHbool = "manager:user:urlAuthBool:"
	REDIS_MANAGER_USER_LOGIN_LOG   = "manager:user:login:log:"
	REDIS_MANAGER_USER_TOKEN       = "manager:user:token:"
	REDIS_MANAGER_USER_SESSION     = "manager:user:session:"
	REDIS_GUILD_SESSION            = "guild:session:"
)

func RedisKeySmsNumIp(ip string) string {
	return fmt.Sprintf("sms:num:ip:%s", ip)
}
func RedisKeySmsNumPhone(phone string) string {
	return fmt.Sprintf("sms:num:phone:%s", phone)
}

func RedisKeySmsNumTime(phone string) string {
	return fmt.Sprintf("sms:num:time:%s", phone)
}
