package service

import (
	"Webook/webook/internal/repository"
	"Webook/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

const codeTplId = "1877556"

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc *sms.Service
}

// Send 发送验证码: biz 业务名称，phone 手机号
func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()

	// 存入 redis
	if err := svc.repo.Store(ctx, biz, phone, code); err != nil {
		return err
	}

	// 发送短信
	if err := svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone); err != nil {
		return err
	}

	return nil
}

// Verify 验证验证码: biz 业务名称，phone 手机号，inputCode 输入的验证码
func (svc *CodeService) Verify(ctx context.Context, biz string,
	phone string, inputCode string) (bool, error) {

	return false, nil
}

func (svc *CodeService) generateCode() string {
	// 生成 6 位随机数 [000000 ~ 999999]
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
