package domain

type Interactive struct {
	// 业务
	Biz   string
	BizId int64

	// 与业务有关
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64

	// 与具体的 userId 有关
	Liked     bool
	Collected bool
}
