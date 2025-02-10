package failover

import (
	"Webook/webook/internal/service/sms"
	"context"
	"log"
	"sync/atomic"
)

type TimeOutFailoverSMSService struct {
	svcs []sms.Service
	// 当前轮询到的服务索引
	idx uint32
	// 连续发送失败次数
	cnt uint32
	// 连续发送失败次数阈值
	threshold uint32
}

func NewTimeOutFailoverSMSService(svcs []sms.Service, threshold uint32) *TimeOutFailoverSMSService {
	return &TimeOutFailoverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}

func (s *TimeOutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	cnt := atomic.LoadUint32(&s.cnt)
	idx := atomic.LoadUint32(&s.idx)

	// 超过阈值，执行切换
	if cnt > s.threshold {
		newIdx := (idx + 1) % uint32(len(s.svcs))
		if atomic.CompareAndSwapUint32(&s.idx, idx, newIdx) {
			// 重置 cnt
			atomic.StoreUint32(&s.cnt, 0)
		}
		idx = newIdx
	}
	svc := s.svcs[int(idx)]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case nil:
		// 发送成功，重置 cnt
		atomic.StoreUint32(&s.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		atomic.AddUint32(&s.cnt, 1)
	default:
		log.Println(err)
	}
	return err
}
