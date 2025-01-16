package service

import (
	"Webook/webook/internal/repository"
	"Webook/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

const codeTplId = "1877556"

var (
	ErrCodeSendTooFrequent    = repository.ErrCodeSendTooFrequent
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CodeServiceStruct struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &CodeServiceStruct{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// Send 发送验证码: biz 业务名称，phone 手机号
func (svc *CodeServiceStruct) Send(ctx context.Context, biz string, phone string) error {
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
func (svc *CodeServiceStruct) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeServiceStruct) generateCode() string {
	// 生成 6 位随机数 [000000 ~ 999999]
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
