package failover

import (
	"Webook/webook/internal/service/sms"
	"context"
	"errors"
	"log"
	"sync/atomic"
)

type FailoverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSService(svcs []sms.Service) *FailoverSMSService {
	return &FailoverSMSService{
		svcs: svcs,
		idx:  0,
	}
}

func (s *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range s.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		// 发送失败，记录日志
		log.Println(err)
	}
	return errors.New("all sms services failed")
}

// idx 是当前轮询到的服务索引
func (s *FailoverSMSService) SendV2(ctx context.Context, tpl string, args []string, numbers ...string) error {
	// 用下一个节点作为当前轮询到的服务索引
	idx := atomic.AddUint64(&s.idx, 1)
	len := uint64(len(s.svcs))
	for i := idx; i < idx+len; i++ {
		svc := s.svcs[i%len]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		default:
			log.Println(err)
		}
	}
	return errors.New("all sms services failed")
}
