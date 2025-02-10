package tencent

import (
	"Webook/webook/pkg/ratelimit"
	"context"
	"fmt"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
}

func NewService(c *sms.Client, appId, signName string, limiter ratelimit.Limiter) *Service {
	return &Service{
		appId:    &appId,
		signName: &signName,
		client:   c,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = &tplId
	req.PhoneNumberSet = str2strPtr(numbers...)
	req.TemplateParamSet = str2strPtr(args...)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}

	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("send sms failed, status code: %s, status message: %s",
				*status.Code, *status.Message)
		}
	}
	return nil
}

func str2strPtr(src ...string) []*string {
	res := make([]*string, len(src))
	for i, s := range src {
		res[i] = &s
	}
	return res
}
