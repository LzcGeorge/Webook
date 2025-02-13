package ioc

import (
	"Webook/webook/internal/service/oauth2/wechat"
	"os"
)

func InitWechatService() wechat.Service {
	appID, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		// panic("找不到环境变量 WECHAT_APP_ID")
		appID = "wx6666666666666666"
	}
	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		// panic("找不到环境变量 WECHAT_APP_SECRET")
		appSecret = "66666666666666666666666666666666"
	}
	return wechat.NewService(appID, appSecret)
}
