package ioc

import (
	"Webook/webook/internal/service/sms"
	"Webook/webook/internal/service/sms/memory"
)

// InitSMSService 初始化短信服务
func InitSMSService() sms.Service {
	// 采用本地内存实现
	return memory.NewService()
}
