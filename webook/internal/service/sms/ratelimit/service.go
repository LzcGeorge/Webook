package ratelimit

import (
	"Webook/webook/internal/service/sms"
	"Webook/webook/pkg/ratelimit"
	"context"
	"fmt"
)

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) *RatelimitSMSService {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

// 采用装饰者模式 Send
func (s *RatelimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 判断是否限流
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return fmt.Errorf("短信服务判断是否限流出现问题,%w", err)
	}
	if limited {
		return fmt.Errorf("短信服务发送短信过于频繁,触发限流")
	}

	err = s.svc.Send(ctx, tplId, args, numbers...)

	// 之后的操作

	return err
}
